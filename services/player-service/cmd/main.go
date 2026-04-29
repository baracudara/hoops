package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/baracudara/hoops/player-service/internal/app"
	"github.com/baracudara/hoops/player-service/internal/config"
)

func main() {
    cfg := config.MustLoad()

    log := slog.New(
        slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
    )

    application := app.New(log, cfg)

    go func() {
        application.GRPCServer.MustRun()
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
    <-stop

    application.GRPCServer.Stop()
    log.Info("player-service stopped")
}