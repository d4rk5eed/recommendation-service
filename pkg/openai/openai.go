package openai

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client представляет клиент для работы с OpenAI API
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	model      string
}

// Message представляет сообщение в чате
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request представляет запрос к OpenAI API
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// Response представляет ответ от OpenAI API
type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Error представляет ошибку OpenAI API
type Error struct {
	StatusCode int
	Body       string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("OpenAI API error (status %d): %s", e.StatusCode, e.Message)
}

// NewClient создает новый клиент OpenAI
func NewClient(apiKey string, opts ...Option) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &Client{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1",
		httpClient: &http.Client{Timeout: 30 * time.Second, Transport: tr},
		model:      "gpt-3.5-turbo",
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Option представляет функцию для настройки клиента
type Option func(*Client)

// WithHTTPClient устанавливает пользовательский HTTP клиент
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL устанавливает базовый URL API
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithModel устанавливает модель по умолчанию
func WithModel(model string) Option {
	return func(c *Client) {
		c.model = model
	}
}

// ChatCompletion выполняет запрос к API чата
func (c *Client) ChatCompletion(ctx context.Context, messages []Message) (string, error) {
	reqBody := Request{
		Model:    c.model,
		Messages: messages,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return "", fmt.Errorf("encode request: %w", err)
	}

	fmt.Printf("%v", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", &buf)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr Response
		if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error.Message != "" {
			return "", &Error{
				StatusCode: resp.StatusCode,
				Body:       string(body),
				Message:    apiErr.Error.Message,
			}
		}
		return "", &Error{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			Message:    "unexpected status code",
		}
	}

	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return "", errors.New("no choices in response")
	}

	return apiResp.Choices[0].Message.Content, nil
}

// SimpleCompletion выполняет упрощенный запрос с одним сообщением пользователя
func (c *Client) SimpleCompletion(ctx context.Context, prompt string) (string, error) {
	return c.ChatCompletion(ctx, []Message{
		{
			Role:    "user",
			Content: prompt,
		},
	})
}
