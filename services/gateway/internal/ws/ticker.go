package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TicketStore struct {
    client *redis.Client
}

func NewTicketStore(client *redis.Client) *TicketStore {
    return &TicketStore{client: client}
}

func (t *TicketStore) Create(ctx context.Context, playerID string) (string, error) {
    ticket := uuid.New().String()

    err := t.client.Set(ctx, "ws:ticket:"+ticket, playerID, 30*time.Second).Err()
    if err != nil {
        return "", fmt.Errorf("failed to create ticket: %w", err)
    }

    return ticket, nil
}

func (t *TicketStore) Consume(ctx context.Context, ticket string) (string, error) {
    key := "ws:ticket:" + ticket

    playerID, err := t.client.Get(ctx, key).Result()
    if err != nil {
        return "", fmt.Errorf("ticket not found or expired: %w", err)
    }

    // удаляем ticket — он одноразовый
    t.client.Del(ctx, key)

    return playerID, nil
}