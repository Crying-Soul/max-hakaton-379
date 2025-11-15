# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Установка рабочей директории
WORKDIR /build

# Копирование go.mod и go.sum для кеширования слоев
COPY go.mod go.sum* ./

# Скачивание зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
# Флаги для уменьшения размера и оптимизации
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /build/maxbot ./cmd/bot/main.go

# Проверка что бинарник создался
RUN ls -lah /build/maxbot

# Stage 2: Runtime
FROM alpine:latest

# Установка зависимостей для runtime (ca-certificates для HTTPS)
RUN apk add --no-cache ca-certificates tzdata

# Создание non-root пользователя для безопасности
RUN addgroup -g 1000 bot && adduser -D -u 1000 -G bot bot

# Установка рабочей директории
WORKDIR /app

# Копирование бинарника из builder стадии
COPY --from=builder /build/maxbot /app/maxbot

# Копирование папки с миграциями
COPY --chown=bot:bot db/migrations /app/db/migrations

# Копирование сертификатов для HTTPS (если они присутствуют)
COPY --chown=bot:bot certs /app/certs

# Копирование .env файла (если существует)
# В production рекомендуется передавать переменные через docker-compose или secrets
#COPY --chown=bot:bot .env* /app/

# Смена пользователя на non-root для безопасности
USER bot

# Здоровье чека (если есть возможность)
# Раскомментируйте если добавите HTTP сервер для хеалт-чека
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD /app/maxbot health

# Запуск приложения
CMD ["/app/maxbot"]
