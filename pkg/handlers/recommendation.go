package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	authstrategy "recommedation-service/pkg/auth_strategy"
	"recommedation-service/pkg/models"
	ai "recommedation-service/pkg/openai"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func HandleRecommendationRequest(w http.ResponseWriter, r *http.Request, conf models.ServiceConfig) {
	var req models.RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var recommendations models.RecommendationResponse
	var err error

	recommendations, err = getLLMRecommendations(req, conf[req.Algorithm])
	if err != nil && err.Error() == "invalid algorithm" {
		respondWithError(w, http.StatusBadRequest, err.Error())
		// return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		// return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"recommendations": recommendations.Recommendations,
		"algorithm":       req.Algorithm,
	})
}

func getLLMRecommendations(req models.RecommendationRequest, conf map[string]string) (models.RecommendationResponse, error) {
	var token string

	var err error
	var response models.RecommendationResponse

	if conf == nil {
		return response, errors.New("invalid algorithm")
	}

	strategy := authstrategy.StrategyFactory(conf)
	token, err = strategy.Execute()
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return response, err
	}
	// Создаем клиент с тестовым URL
	client := ai.NewClient(token, ai.WithBaseURL(conf["url"]), ai.WithModel(conf["model"]))

	// Вызываем метод
	resp, err := client.GetRecommendation(context.Background(), req)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
		return response, err
	}

	// Декодируем JSON
	err = json.Unmarshal([]byte(cleanJSONString(resp)), &response)
	if err != nil {
		log.Fatalf("Ошибка декодирования JSON: %v", err)
		return response, err
	}

	return response, nil
}

func getRAGRecommendations(_ models.RecommendationRequest) ([]models.Recommendation, error) {
	return nil, errors.New("not implemented")
}

func cleanJSONString(jsonStr string) string {
	// Удаляем маркер ```json в начале, если есть
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	// Удаляем апострофы ``` в начале/конце, если есть
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	// Удаляем лишние пробелы и переносы
	return strings.TrimSpace(jsonStr)
}
