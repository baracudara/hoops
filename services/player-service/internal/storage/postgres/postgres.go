package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baracudara/hoops/player-service/internal/config"
	"github.com/baracudara/hoops/player-service/internal/domain/dto"
	"github.com/baracudara/hoops/player-service/internal/domain/models"
	"github.com/baracudara/hoops/player-service/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
    psql *pgxpool.Pool
}

func New(cfg *config.Postgres) (*Storage, error) {
    const op = "storage.postgres.new"

    connString := fmt.Sprintf(
        "postgresql://%s:%s@%s:%d/%s",
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
    )

    pgxConfig, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    pgxConfig.MaxConns = cfg.MaxConns
    pgxConfig.MinConns = cfg.MinConns
    pgxConfig.MaxConnIdleTime = 30 * time.Minute
    pgxConfig.MaxConnLifetime = 1 * time.Hour

    pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{psql: pool}, nil
}



func (s *Storage) SavePlayer(ctx context.Context, player models.Player) (models.Player, error) {
    const op = "storage.postgres.SavePlayer"

    query := `
        INSERT INTO players (id, name, nickname, position, age)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, name, nickname, position, age, games_played, wins, losses, created_at
    `

    var res models.Player
    err := s.psql.QueryRow(ctx, query,
        player.ID,
        player.Name,
        player.Nickname,
        player.Position,
        player.Age,
    ).Scan(
        &res.ID,
        &res.Name,
        &res.Nickname,
        &res.Position,
        &res.Age,
        &res.Stats.GamesPlayed,
        &res.Stats.Wins,
        &res.Stats.Losses,
        &res.CreatedAt,
    )

    if err != nil {
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (s *Storage) GetPlayer(ctx context.Context, uuid string) (models.Player, error) {
    const op = "storage.postgres.GetPlayer"

    query := `
        SELECT id, name, nickname, position, age, games_played, wins, losses, created_at
        FROM players
        WHERE id = $1
    `

    var res models.Player
    err := s.psql.QueryRow(ctx, query, uuid).Scan(
        &res.ID,
        &res.Name,
        &res.Nickname,
        &res.Position,
        &res.Age,
        &res.Stats.GamesPlayed,
        &res.Stats.Wins,
        &res.Stats.Losses,
        &res.CreatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return models.Player{}, fmt.Errorf("%s: %w", op, storage.ErrPlayerNotFound)
        }
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (s *Storage) UpdatePlayer(ctx context.Context, uuid string, dto dto.UpdatePlayer) (models.Player, error) {
    const op = "storage.postgres.UpdatePlayer"

    query := `
        UPDATE players
        SET name = $1, nickname = $2, position = $3, age = $4, updated_at = NOW()
        WHERE id = $5
        RETURNING id, name, nickname, position, age, games_played, wins, losses, created_at
    `

    var res models.Player
    err := s.psql.QueryRow(ctx, query,
        dto.Name,
        dto.Nickname,
        dto.Position,
        dto.Age,
        uuid,
    ).Scan(
        &res.ID,
        &res.Name,
        &res.Nickname,
        &res.Position,
        &res.Age,
        &res.Stats.GamesPlayed,
        &res.Stats.Wins,
        &res.Stats.Losses,
        &res.CreatedAt,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return models.Player{}, fmt.Errorf("%s: %w", op, storage.ErrPlayerNotFound)
        }
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}