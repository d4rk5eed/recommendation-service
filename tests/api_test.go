package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"recommedation-service/pkg/handlers"
	"recommedation-service/pkg/models"
)

func TestRecommendationEndpoint(t *testing.T) {
	tt := []struct {
		name       string
		payload    interface{}
		statusCode int
	}{
		{
			name: "valid llm request",
			payload: models.RecommendationRequest{
				Metrics:   map[string]float64{"cpu": 95},
				Algorithm: "mock-llm",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "invalid algorithm",
			payload: map[string]interface{}{
				"metrics":   map[string]float64{"cpu": 95},
				"algorithm": "invalid",
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/v1/recommendations", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			conf, _ := models.ReadConfig("../config/test.yaml")

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.HandleRecommendationRequest(w, r, conf)
			})
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.statusCode)
			}
		})
	}
}

func TestMockLLM(t *testing.T) {
	// Тест мокового LLM сервера
	body := []byte(`{
		"messages": [{"content": "CPU usage is high"}]
	}`)

	resp, err := http.Post("http://localhost:8081/v1/chat/completions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.MockLLMResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Choices) == 0 {
		t.Error("empty response from mock LLM")
	}
}
