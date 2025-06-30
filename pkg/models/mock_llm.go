package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MockLLMRequest struct {
	Messages []struct {
		Content string `json:"content"`
	} `json:"messages"`
}

type MockLLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func StartMockLLMServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		var req MockLLMRequest
		json.NewDecoder(r.Body).Decode(&req)

		response := MockLLMResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{
					Message: struct {
						Content string `json:"content"`
					}{
						Content: `{"recommendations": [{
							"priority": "high",
							"action": "Check system metrics",
							"details": "Mock recommendation from LLM"
						}]}`,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Mock LLM server running on :8081")
	http.ListenAndServe(":8081", mux)
}
