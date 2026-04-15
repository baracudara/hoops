FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/
RUN go mod download

COPY . .

RUN go build -o auth-server ./services/auth-service/cmd/auth

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-server .

EXPOSE 44044

CMD ["./auth-server"]