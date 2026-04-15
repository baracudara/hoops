FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем корневой go.mod для replace directive
COPY go.mod go.sum ./

# Копируем auth-service
COPY services/auth-service ./services/auth-service

# Билдим из папки auth-service
WORKDIR /app/services/auth-service
RUN go mod download
RUN go build -o /app/auth-server ./cmd/auth

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-server .

EXPOSE 44044

CMD ["./auth-server"]