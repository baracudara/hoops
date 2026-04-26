package authclient

import (
	"context"
	"fmt"

	"github.com/baracudara/hoops/gateway/internal/config"
	"github.com/baracudara/hoops/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
    api auth.AuthClient
}

func New(cfg *config.AuthGRPC) (*Client, error) {
    const op = "clients.auth.New"

    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

    conn, err := grpc.NewClient(addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Client{
        api: auth.NewAuthClient(conn),
    }, nil
}


func (c *Client) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
    const op = "clients.auth.Register"

    res, err := c.api.Register(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (c *Client) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
    const op = "clients.auth.Login"

    res, err := c.api.Login(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (c *Client) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
    const op = "clients.auth.Logout"

    res, err := c.api.Logout(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (c *Client) VerifyAccessToken(ctx context.Context, req *auth.VerifyAccessTokenRequest) (*auth.VerifyAccessTokenResponse, error) {
    const op = "clients.auth.VerifyAccessToken"

    res, err := c.api.VerifyAccessToken(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}


func (c *Client) Refresh(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
    const op = "clients.auth.Refresh"

    res, err := c.api.Refresh(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}