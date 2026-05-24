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
	"github.com/baracudara/hoops/gateway/internal/ws"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
    cfg := config.MustLoad()

    log := slog.New(
        slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
    )

    // gRPC клиент для auth-service
    authClient, err := authclient.New(&cfg.AuthGRPC)
    if err != nil {
        log.Error("failed to connect to auth-service", "err", err)
        os.Exit(1)
    }

    // Redis клиент для tickets
    redisClient := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })

    // WebSocket
    hub := ws.NewHub()
    go hub.Run()

    ticketStore := ws.NewTicketStore(redisClient)
    wsHandler := ws.NewHandler(hub, log, ticketStore)

    // хендлеры
    authHandler := authhandler.New(authClient, cfg.HTTP.RefreshTokenTTL, cfg.HTTP.CookieDomain)

    // роутер
    r := chi.NewRouter()
    r.Use(chimiddleware.Logger)
    r.Use(chimiddleware.Recoverer)

    // открытые роуты
    r.Post("/auth/register", authHandler.Register)
    r.Post("/auth/login", authHandler.Login)
    r.Post("/auth/logout", authHandler.Logout)
    r.Post("/auth/refresh", authHandler.Refresh)



    // защищённые роуты
    r.Group(func(r chi.Router) {
        r.Use(custommiddleware.AuthMiddleware(authClient))
        r.Get("/users/me", authHandler.Me)
        r.Post("/ws/ticket", wsHandler.CreateTicket)
    }) 

    // WebSocket роут
    r.Get("/ws", wsHandler.ServeWS)
	r.Post("/ws/test", func(w http.ResponseWriter, r *http.Request) {
		hub.Broadcast([]byte(`{"type":"test","message":"hello from server"}`))
		w.WriteHeader(http.StatusOK)
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