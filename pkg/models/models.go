package models

import (
	"time"
)

type RecommendationRequest struct {
	Metrics            map[string]float64 `json:"metrics"`
	ProblemClass       string             `json:"problem_class,omitempty"`
	ProblemDescription string             `json:"problem_description,omitempty"`
	Algorithm          string             `json:"algorithm"`
}

type Recommendation struct {
	Priority string   `json:"priority"`
	Action   string   `json:"action"`
	Details  string   `json:"details"`
	Links    []string `json:"links,omitempty"`
}

// Структура для разбора JSON ответа
type RecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
}

type ProblemClass struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Rules       []MetricRule `json:"rules,omitempty"`
}

type MetricRule struct {
	Metric    string  `json:"metric"`
	Condition string  `json:"condition"`
	Weight    float64 `json:"weight"`
}

type Solution struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Steps     []string  `json:"steps"`
	CreatedAt time.Time `json:"created_at"`
}

// OAuthResponse представляет структуру ответа от OAuth сервера
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

var (
	problemClasses map[string]ProblemClass
	solutions      map[string]Solution
)
