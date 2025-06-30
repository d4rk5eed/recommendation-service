# Архитектурная концепция
Сервис будет представлять собой микросервис на Go, который:
- Принимает метрики системы и описание проблемы
- Обрабатывает их с помощью выбранного алгоритма
- Возвращает рекомендации по устранению проблемы

# Компоненты:
API - входная точка для всех запросов
Request Processor - обработчик запросов и маршрутизатор
LLM Proxy - интеграция с облачными или локальным  LLM (OpenAI, Gemini, DeepSeek, Gigachat, YandexGPT и др.)

## Планируемые:
Local RAG Engine - локальная система на основе RAG (Retrieval-Augmented Generation)
Metrics Classifier - классификатор проблем на основе метрик
Vector Knowledge Base - хранилище документации и решений

# Дальнейшее развитие
1. Интеграция с облачными LLM:
  - Добавить поддержку OpenAI, Gemini, DeepSeek, Gigachat, YandexGPT (на выбор)
  - Реализовать кэширование ответов

2. Локальный RAG Engine:
  - Реализовать индексацию внутренней документации
  - Добавить векторное хранилище

3. Классификация проблем:
  - Разработать модель машинного обучения для автоматической классификации
  - Реализовать pipeline для обучения с учителем

4. Мониторинг и логирование:
  - MLflow
  - MLOps

```
curl -X POST http://localhost:8080/v1/recommendations \
-H "Content-Type: application/json" \
-d '{
  "metrics": {
    "cpu_usage": 95.2,
    "memory_usage": 87.5
  },
  "problem_class": "performance_issue",
  "algorithm": "mock-llm"
}'
```
