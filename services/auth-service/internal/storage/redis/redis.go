package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/baracudara/hoops/auth-service/internal/config"
	"github.com/baracudara/hoops/auth-service/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(cfg *config.Redis) (*Redis, error) {

	const op = "storage.redis.new"


	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
        Password: cfg.Password,
        DB:       cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	return &Redis{client: client}, nil
}


func (r *Redis) SaveToken(ctx context.Context, uuid string, token string, ttl time.Duration) error {

	const op = "storage.redis.token.save"

	err := r.client.Set(ctx, "token:"+token, uuid, ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Redis) IsTokenValid(ctx context.Context, token string) (bool, error) {
	const op = "storage.redis.token.valid"

	rslt, err := r.client.Exists(ctx, "token:"+token).Result()

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return rslt > 0, nil
}

func (r *Redis) DeleteToken(ctx context.Context, token string) error {
    const op = "storage.redis.token.delete"

    res, err := r.client.Del(ctx, "token:"+token).Result()
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

	if res == 0 {
        return fmt.Errorf("%s: %w", op, storage.ErrTokenNotFound)
    }

    return nil
}
