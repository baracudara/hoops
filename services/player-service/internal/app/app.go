package app

import (
	"log/slog"

	grpcapp "github.com/baracudara/hoops/player-service/internal/app/grpc"
	"github.com/baracudara/hoops/player-service/internal/config"
	"github.com/baracudara/hoops/player-service/internal/services/player"
	"github.com/baracudara/hoops/player-service/internal/storage/postgres"
)

type App struct {
    GRPCServer *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
    storage, err := postgres.New(&cfg.Postgres)
    if err != nil {
        panic(err)
    }

    playerService := player.New(log, storage, storage, storage)

    grpcApp := grpcapp.New(log, playerService, cfg.GRPC.Port)

    return &App{
        GRPCServer: grpcApp,
    }
}