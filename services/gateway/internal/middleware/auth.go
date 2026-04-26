package middleware

import (
	"context"
	"net/http"
	"strings"

	authclient "github.com/baracudara/hoops/gateway/internal/clients/auth"
	"github.com/baracudara/hoops/protos/gen/go/auth"
)

type contextKey string

const UserKey contextKey = "user"

func AuthMiddleware(client *authclient.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "invalid token format", http.StatusUnauthorized)
                return
            }

            token := parts[1]

            res, err := client.VerifyAccessToken(r.Context(), &auth.VerifyAccessTokenRequest{
                AccessToken: token,
            })
            if err != nil || !res.Valid {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            // кладём данные пользователя в контекст
            ctx := context.WithValue(r.Context(), UserKey, res)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}