package main

import (
	"log"

	"recommedation-service/pkg/models"
)

func main() {
	// Запуск мокового LLM сервера
	models.StartMockLLMServer()

	log.Println("Starting main server on :8081")
}
