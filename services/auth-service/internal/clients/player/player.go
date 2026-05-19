package playerclient

import (
	"context"
	"fmt"

	"github.com/baracudara/hoops/auth-service/internal/config"
	"github.com/baracudara/hoops/protos/gen/go/player"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
    api player.PlayerClient
}

func New(cfg *config.PlayerGRPC) (*Client, error) {
    const op = "clients.player.New"

    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

    var opts []grpc.DialOption

    if cfg.Secure {
        creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
        if err != nil {
            return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
        }
        opts = append(opts, grpc.WithTransportCredentials(creds))
    } else {
        opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    }

    conn, err := grpc.NewClient(addr, opts...)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Client{
        api: player.NewPlayerClient(conn),
    }, nil
}

func (c *Client) CreatePlayer(ctx context.Context, req *player.CreatePlayerRequest) (*player.CreatePlayerResponse, error) {
    const op = "clients.player.CreatePlayer"

    res, err := c.api.CreatePlayer(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}