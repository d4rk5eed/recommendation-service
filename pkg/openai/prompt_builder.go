package openai

import (
	"context"
	"fmt"
	"strings"

	"recommedation-service/pkg/models"
)

// BuildPromptFromRequest преобразует RecommendationRequest в промпт для OpenAI
func BuildPromptFromRequest(req models.RecommendationRequest) string {
	var promptParts []string

	// Добавляем метрики
	if len(req.Metrics) > 0 {
		metricsPart := "Системные метрики:\n"
		for name, value := range req.Metrics {
			metricsPart += fmt.Sprintf("- %s: %.2f\n", name, value)
		}
		promptParts = append(promptParts, metricsPart)
	}

	// Добавляем класс проблемы
	if req.ProblemClass != "" {
		promptParts = append(promptParts, fmt.Sprintf("Класс проблемы: %s", req.ProblemClass))
	}

	// Добавляем описание проблемы
	if req.ProblemDescription != "" {
		promptParts = append(promptParts, fmt.Sprintf("Описание проблемы: %s", req.ProblemDescription))
	}

	// Добавляем инструкцию
	instruction := `Проанализируй проблему и предоставь рекомендации по устранению.
Ответ должен быть строго структурированным в формате json и содержать только ключ recommendations - список объектов,
	каждый объект содержит
1. priority: Приоритет проблемы (high, medium, low)
2. action: Конкретные действия для решения на русском языке
3. details: Детализацию рекомендаций на русском языке
4. links: Ссылки на документацию (если уместно)

Других объектов быть не должно.
Ответ должен быть готовым к декодированию с помощью Unmarshal. markdown разметку нужно удалить`

	promptParts = append(promptParts, instruction)

	return strings.Join(promptParts, "\n\n")
}

// GetRecommendation получает рекомендации на основе RecommendationRequest
func (c *Client) GetRecommendation(ctx context.Context, req models.RecommendationRequest) (string, error) {
	prompt := BuildPromptFromRequest(req)
	return c.SimpleCompletion(ctx, prompt)
}
