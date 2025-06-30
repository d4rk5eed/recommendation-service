# README для сервиса рекомендаций по техническому обслуживанию

## Описание сервиса

Сервис рекомендаций по устранению технических проблем на основе метрик системы. Поддерживает два режима работы:
- С реальным LLM API (GigaChat)
- С моковым LLM сервером для тестирования

## Требования

- Go 1.24+
- Docker (опционально)
- GitLab Runner (для CI/CD)

## Установка и запуск

### Локальная разработка

1. Клонируйте репозиторий:
```bash
git clone git@gitlab.cti.ru:devops/recommendation-service.git
cd recommendation-service
```


3. Соберите и запустите сервис:
```bash
make run
```

### Запуск с моковым LLM

Для тестирования без реального API:
```bash
make build-mock-llm
make run-mock-llm
```

## Конфигурация

Создайте файл конфигурации `config/dev.yaml`:
```yaml
api:
  url: "https://api.example.com"
  key: "your-api-key"
```

Или используйте переменные окружения:
```bash
export API_KEY="your-api-key"
export GIGACHAT_API_PERS="your-scope" //только для Oauth сервиса Gigachat
```

## Docker

### Сборка образа
```bash
make docker-build
```

### Запуск контейнера
```bash
API_KEY="<some key> "GIGACHAT_API_PERS="<some key>" make docker-run
```

Или с кастомными параметрами:
```bash
docker run -p 8080:8080 -e API_KEY="your-key" -e GIGACHAT_API_PERS="your-scope" recommendation-service:latest
```

## Тестирование

### Юнит-тесты
```bash
make test
```

### Все тесты
```bash
make test-all
```

## Доступные команды Makefile

| Команда              | Описание                                      |
|----------------------|-----------------------------------------------|
| `make build`         | Собрать бинарный файл                         |
| `make run`           | Запустить сервер локально                     |
| `make test`          | Запустить юнит-тесты                          |
| `make test-integration` | Запустить интеграционные тесты              |
| `make docker-build`  | Собрать Docker образ                          |
| `make docker-run`    | Запустить в Docker                            |
| `make fmt`           | Форматировать код                             |
| `make lint`          | Запустить линтер                              |
| `make clean`         | Очистить проект                               |

## Переменные окружения

| Переменная           | Описание                     | По умолчанию      |
|----------------------|-----------------------------|-------------------|
| `API_KEY`            | Ключ для доступа к API       | -                 |
| `GIGACHAT_API_PERS`  | Scope для OAuth              | -                 |
| `MAIN_PORT`          | Порт основного сервера       | `8080`            |
| `MOCK_LLM_PORT`      | Порт мокового LLM            | `8081`            |


## Другие ссылки
- [Описание архитектуры](docs/arch.md)
- [Описание API](docs/api.md)
