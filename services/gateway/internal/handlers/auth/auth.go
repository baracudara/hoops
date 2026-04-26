package authhandler

import (
	"encoding/json"
	"net/http"
	"time"

	"io"

	authclient "github.com/baracudara/hoops/gateway/internal/clients/auth"
	"github.com/baracudara/hoops/gateway/internal/middleware"
	"github.com/baracudara/hoops/protos/gen/go/auth"
	"google.golang.org/protobuf/encoding/protojson"
)

type Handler struct {
    client          *authclient.Client
    refreshTokenTTL time.Duration
    cookieDomain    string
}

func New(client *authclient.Client, refreshTokenTTL time.Duration, cookieDomain string) *Handler {
    return &Handler{
        client:          client,
        refreshTokenTTL: refreshTokenTTL,
        cookieDomain:    cookieDomain,
    }
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    // стало
    var req auth.RegisterRequest
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    if err := protojson.Unmarshal(body, &req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    res, err := h.client.Register(r.Context(), &req)
    if err != nil {
        http.Error(w, "failed to register", http.StatusInternalServerError)
        return
    }

    // refresh token кладём в cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "refreshToken",
        Value:    res.RefreshToken,
        HttpOnly: true,
        Domain:   h.cookieDomain,
        Path:     "/",
        Expires:  time.Now().Add(h.refreshTokenTTL),
    })

    // access token возвращаем в JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": res.AccessToken,
    })
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req auth.LoginRequest
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    if err := protojson.Unmarshal(body, &req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    res, err := h.client.Login(r.Context(), &req)

    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return 
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "refreshToken",
        Value:    res.RefreshToken,
        HttpOnly: true,
        Domain:   h.cookieDomain,
        Path:     "/",
        Expires:  time.Now().Add(h.refreshTokenTTL),
    })


    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": res.AccessToken,
    })


}



func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

    cookie, err := r.Cookie("refreshToken")
    if err != nil {
        http.Error(w, "refresh token not found", http.StatusUnauthorized)
        return
    }

    _, err = h.client.Logout(r.Context(), &auth.LogoutRequest{
        Token: cookie.Value,
    })

    if err != nil {
        http.Error(w, "failed to logout", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "refreshToken",
        Value:   "",
        Expires: time.Unix(0, 0),
        MaxAge:  -1,
        Path:    "/",
    })

    w.WriteHeader(http.StatusOK)

}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
    // достаём данные пользователя из контекста
    // которые положил туда middleware
    user, ok := r.Context().Value(middleware.UserKey).(*auth.VerifyAccessTokenResponse)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // возвращаем в JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "uuid": user.Uuid,
        "role": user.Role,
    })
}



func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
    // достаём refresh token из cookie
    cookie, err := r.Cookie("refreshToken")
    if err != nil {
        http.Error(w, "refresh token not found", http.StatusUnauthorized)
        return
    }

    res, err := h.client.Refresh(r.Context(), &auth.RefreshRequest{
        RefreshToken: cookie.Value,
    })
    if err != nil {
        http.Error(w, "failed to refresh tokens", http.StatusUnauthorized)
        return
    }

    // обновляем cookie с новым refresh token
    http.SetCookie(w, &http.Cookie{
        Name:     "refreshToken",
        Value:    res.RefreshToken,
        HttpOnly: true,
        Domain:   h.cookieDomain,
        Path:     "/",
        Expires:  time.Now().Add(h.refreshTokenTTL),
    })

    // возвращаем новый access token
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": res.AccessToken,
    })
}