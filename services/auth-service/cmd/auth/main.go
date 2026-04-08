package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/baracudara/hoops/auth-service/internal/app"
	"github.com/baracudara/hoops/auth-service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd      = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)


	application := app.New(log, cfg)

	go func () {
		application.GRPSServer.MustRun()
	} ()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop





	
	application.GRPSServer.Stop()
	log.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
    switch env {
    case envLocal:
        return slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envDev:
        return slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case envProd:
        return slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        )
    default:
        return slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    }
}