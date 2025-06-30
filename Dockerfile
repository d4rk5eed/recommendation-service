# Этап сборки
FROM golang:1.24-alpine AS base

RUN apk add --no-cache make


WORKDIR /app
COPY . .
RUN make build

# Этап запуска
FROM alpine:latest

WORKDIR /app
RUN ls

COPY --from=base /app/bin/recommendation-service /app/recommendation-service
COPY --from=base /app/config/ /app/config/

# Переменные среды
ENV API_KEY=${API_KEY} \
    GIGACHAT_API_PERS=${GIGACHAT_API_PERS}
RUN echo ${GIGACHAT_API_PERS}
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/health || exit 1

# ENTRYPOINT ["/app/recommendation-service"]
CMD ["/app/recommendation-service", "--config", "/app/config/prod.yaml"]
