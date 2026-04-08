package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baracudara/hoops/auth-service/internal/config"
	"github.com/baracudara/hoops/auth-service/internal/domain/dto"
	"github.com/baracudara/hoops/auth-service/internal/domain/models"
	"github.com/baracudara/hoops/auth-service/internal/storage"
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

	pgxConifg, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pgxConifg.MaxConns = cfg.MaxConns
	pgxConifg.MinConns = cfg.MinConns
	pgxConifg.MaxConnIdleTime = 30 * time.Minute
	pgxConifg.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConifg)

	if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
	
    return &Storage{psql: pool}, nil
}


func (s *Storage) SaveUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.postgres.save"

    query := `
        INSERT INTO users (id, name, nickname, email, phone, google_id, pass_hash, role, trust_rating)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, name, nickname, email, phone, google_id, role, trust_rating, created_at
    `


	var res models.User
    err := s.psql.QueryRow(ctx, query,
        user.ID,
        user.Name,
        user.Nickname,
        user.Email,
        user.Phone,
        user.GoogleID,
        user.PassHash,
        user.Role,
        user.TrustRating,
    ).Scan(
        &res.ID,
        &res.Name,
        &res.Nickname,
        &res.Email,
        &res.Phone,
        &res.GoogleID,
        &res.Role,
        &res.TrustRating,
        &res.CreatedAt,
    )

    if err != nil {
        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}


func (s *Storage) GetUser(ctx context.Context, dto dto.Login) (models.User, error) {

    const op = "storage.postgres.user.get"

    query := `
        SELECT id, name, nickname, email, phone, google_id, pass_hash, role, trust_rating, created_at  
        FROM users
        WHERE email = $1 OR phone = $2 OR google_id = $3
        LIMIT 1
    `

    var user models.User

    err := s.psql.QueryRow(ctx, query, dto.Email, dto.Phone, dto.GoogleID).Scan(
        &user.ID,
        &user.Name,
        &user.Nickname,
        &user.Email,
        &user.Phone,
        &user.GoogleID,
        &user.PassHash,  
        &user.Role,
        &user.TrustRating,
        &user.CreatedAt,

    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
        }

        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }


    return user, nil
}
