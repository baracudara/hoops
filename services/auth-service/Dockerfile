# Стадия 1 — сборка
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-server ./cmd/auth

# Стадия 2 — финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-server .

EXPOSE 44044

CMD ["./auth-server"]