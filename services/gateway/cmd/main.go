package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authclient "github.com/baracudara/hoops/gateway/internal/clients/auth"
	"github.com/baracudara/hoops/gateway/internal/config"
	authhandler "github.com/baracudara/hoops/gateway/internal/handlers/auth"
	custommiddleware "github.com/baracudara/hoops/gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	authClient, err := authclient.New(&cfg.AuthGRPC)

	if err  != nil {
		log.Error("failed to connect to auth-service", "err", err)
		os.Exit(1)
	}

	authHandler := authhandler.New(authClient, cfg.HTTP.RefreshTokenTTL, cfg.HTTP.CookieDomain)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Post("/auth/register", authHandler.Register)
    r.Post("/auth/login", authHandler.Login)
    r.Post("/auth/logout", authHandler.Logout)
	r.Post("/auth/refresh", authHandler.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.AuthMiddleware(authClient))
		r.Get("/users/me", authHandler.Me)
	})

	


    srv := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
        Handler: r,
    }

	log.Info("gateway started", "port", cfg.HTTP.Port)

	go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Error("server error", "err", err)
        }
    }()


	stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
    <-stop

    log.Info("gateway stopped")


}