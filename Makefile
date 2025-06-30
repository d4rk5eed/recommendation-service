# Makefile для сервиса рекомендаций по техническому обслуживанию

# Конфигурация
BINARY_NAME = bin/recommendation-service
MOCK_LLM_BINARY_NAME = bin/mock-llm
MOCK_LLM_PORT = 8081
MAIN_PORT = 8080
DOCKER_IMAGE = recommendation-service:latest

.PHONY: all build run test test-integration docker-build docker-run clean

all: build

# Сборка проекта
build:
	go build -o ${BINARY_NAME} cmd/recommendation-service/recommendation-service.go

# Запуск основного сервера и мокового LLM
run: build
	@echo "Запуск основного сервера на порту ${MAIN_PORT}"
	./${BINARY_NAME} --config config/dev.yaml
	@sleep 1 # Даем время для запуска
	@echo "Серверы запущены. Используйте Ctrl+C для остановки"

	# Сборка проекта
build-mock-llm:
		go build -o ${MOCK_LLM_BINARY_NAME} cmd/mock-llm/mock-llm.go

# Запуск основного сервера и мокового LLM
run-mock-llm: build-mock-llm
		@echo "Запуск мокового LLM на порту ${MOCK_LLM_PORT}"
		./${MOCK_LLM_BINARY_NAME} &
		@sleep 1 # Даем время для запуска
		@echo "Серверы запущены. Используйте Ctrl+C для остановки"

# Запуск всех тестов
test-all:
	go test -v ./tests/

# Запуск юнит-тестов
test:
	go test -v ./tests/ -run TestRecommendationEndpoint

# Запуск интеграционных тестов (требует запущенных серверов)
test-integration:
	@echo "Запуск интеграционных тестов..."
	go test -v ./tests/ -run TestMockLLM

# Сборка Docker образа
docker-build:
	docker build -t ${DOCKER_IMAGE} .

# Запуск в Docker
docker-run:
	@echo "Запуск сервиса в Docker..."
	# docker run -p ${MAIN_PORT}:${MAIN_PORT} -p ${MOCK_LLM_PORT}:${MOCK_LLM_PORT} ${DOCKER_IMAGE}
	docker run  --rm \
  		-p ${MAIN_PORT}:${MAIN_PORT} \
		-e "CONFIG_PATH=${CONFIG_PATH}" \
		-e "API_KEY=${API_KEY}" \
		-e "GIGACHAT_API_PERS=${GIGACHAT_API_PERS}"
  		--name recommendations \
  		${DOCKER_IMAGE}

# Очистка
clean:
	rm -f ${BINARY_NAME}
	go clean

# Форматирование кода
fmt:
	gofmt -w .

# Запуск линтера
lint:
	golangci-lint run

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build         - Собрать бинарный файл"
	@echo "  run           - Запустить сервер локально"
	@echo "  test          - Запустить юнит-тесты"
	@echo "  test-integration - Запустить интеграционные тесты (требует запущенных серверов)"
	@echo "  docker-build  - Собрать Docker образ"
	@echo "  docker-run    - Запустить в Docker"
	@echo "  fmt           - Форматировать код"
	@echo "  lint          - Запустить линтер"
	@echo "  clean         - Очистить проект"
