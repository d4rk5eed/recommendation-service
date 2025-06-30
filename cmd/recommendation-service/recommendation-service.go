package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"recommedation-service/pkg/handlers"
	"recommedation-service/pkg/models"

	"github.com/gorilla/mux"
)

func main() {
	configPath := flag.String("config", "config/test.yaml", "Path to configuration file")
	flag.Parse()

	r := mux.NewRouter()
	conf, err := models.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Регистрация хендлеров
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/v1/recommendations", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRecommendationRequest(w, r, conf)
	}).Methods("POST")
	// r.HandleFunc("/v1/problem-classes", createProblemClass).Methods("POST")
	// r.HandleFunc("/v1/training", trainClassifier).Methods("POST")

	port := fmt.Sprintf(":%s", models.GetEnvWithDefault("PORT", "8090"))
	log.Printf("Starting main server on %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
