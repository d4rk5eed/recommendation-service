package authstrategy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"recommedation-service/pkg/models"
)

type GigachatStrategy struct {
	Conf map[string]string
}

func (s GigachatStrategy) Execute() (string, error) {
	auth, err := GetOAuthToken("https://ngw.devices.sberbank.ru:9443/api/v2/oauth", s.Conf["key"], "GIGACHAT_API_PERS", os.Getenv("GIGACHAT_API_PERS"))
	if err != nil {
		return "", fmt.Errorf("error getting OAuth token: %w", err)
	}
	return auth.AccessToken, nil
}

// GetOAuthToken выполняет запрос для получения OAuth токена
func GetOAuthToken(apiURL, auth, scope, rqUID string) (*models.OAuthResponse, error) {
	// Формируем данные для запроса
	data := url.Values{}
	data.Set("scope", scope)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tr,
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", rqUID)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("oauth server returned status %d: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var oauthResp models.OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &oauthResp, nil
}
