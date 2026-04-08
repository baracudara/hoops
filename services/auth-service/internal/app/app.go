package app

import (
	"log/slog"

	grpcapp "github.com/baracudara/hoops/auth-service/internal/app/grpc"
	"github.com/baracudara/hoops/auth-service/internal/config"
	"github.com/baracudara/hoops/auth-service/internal/services/auth"
	"github.com/baracudara/hoops/auth-service/internal/storage/postgres"
	"github.com/baracudara/hoops/auth-service/internal/storage/redis"
)


type App struct {
	GRPSServer *grpcapp.App

}

func New(
	log *slog.Logger, 
	cfg  *config.Config, 

	
	) *App {
		storage, err := postgres.New(&cfg .Postgres)

		if err != nil {
			panic(err)
		}

		cache, err := redis.New(&cfg .Redis)

		if err != nil {
			panic(err)
		}

		authService := auth.New(log, storage, storage, cache, cache, cache, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL, cfg.JWT.Secret)

		grpcApp := grpcapp.New(log, authService, cfg.GRPC.Port)

		return &App{
			GRPSServer: grpcApp,

		}
}
