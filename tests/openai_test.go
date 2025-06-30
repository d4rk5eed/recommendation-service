package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ai "recommedation-service/pkg/openai"
)

func TestClient_ChatCompletion(t *testing.T) {
	// Тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		response := `{
			"choices": [{
				"message": {
					"content": "Test response"
				}
			}]
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer ts.Close()

	// Создаем клиент с тестовым URL
	client := ai.NewClient("test-key", ai.WithBaseURL(ts.URL))

	// Вызываем метод
	resp, err := client.ChatCompletion(context.Background(), []ai.Message{
		{
			Role:    "user",
			Content: "Test message",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "Test response" {
		t.Errorf("expected 'Test response', got '%s'", resp)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	// Тестовый сервер с ошибкой
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		response := `{
			"error": {
				"message": "Invalid request"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer ts.Close()

	client := ai.NewClient("test-key", ai.WithBaseURL(ts.URL))

	_, err := client.ChatCompletion(context.Background(), []ai.Message{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*ai.Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}

	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestClient_Timeout(t *testing.T) {
	// Тестовый сервер с задержкой
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte(`{"choices": [{"message": {"content": "ok"}}]}`))
	}))
	defer ts.Close()

	// Клиент с маленьким таймаутом
	httpClient := &http.Client{Timeout: 50 * time.Millisecond}
	client := ai.NewClient("test-key", ai.WithBaseURL(ts.URL), ai.WithHTTPClient(httpClient))

	_, err := client.ChatCompletion(context.Background(), []ai.Message{})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}
