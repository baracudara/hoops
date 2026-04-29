package player

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/baracudara/hoops/player-service/internal/domain/dto"
	"github.com/baracudara/hoops/player-service/internal/domain/models"
	"github.com/baracudara/hoops/player-service/internal/lib/logger/sl"
	"github.com/google/uuid"
)

type Player struct {
    log             *slog.Logger
    playerSaver     PlayerSaver
    playerProvider  PlayerProvider
    playerUpdater   PlayerUpdater
}

type PlayerSaver interface {
    SavePlayer(ctx context.Context, player models.Player) (models.Player, error)
}

type PlayerProvider interface {
    GetPlayer(ctx context.Context, uuid string) (models.Player, error)
}

type PlayerUpdater interface {
    UpdatePlayer(ctx context.Context, uuid string, dto dto.UpdatePlayer) (models.Player, error)
}

func New(
    log *slog.Logger,
    playerSaver PlayerSaver,
    playerProvider PlayerProvider,
    playerUpdater PlayerUpdater,
) *Player {
    return &Player{
        log:            log,
        playerSaver:    playerSaver,
        playerProvider: playerProvider,
        playerUpdater:  playerUpdater,
    }
}

func (p *Player) CreatePlayer(ctx context.Context, player models.Player) (models.Player, error) {
    const op = "services.player.CreatePlayer"

    log := p.log.With(slog.String("op", op))
    log.Info("creating player")

    player.ID = uuid.New().String()

    res, err := p.playerSaver.SavePlayer(ctx, player)
    if err != nil {
        log.Error("failed to save player", sl.Err(err))
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (p *Player) GetPlayer(ctx context.Context, uuid string) (models.Player, error) {
    const op = "services.player.GetPlayer"

    log := p.log.With(slog.String("op", op))
    log.Info("getting player")

    res, err := p.playerProvider.GetPlayer(ctx, uuid)
    if err != nil {
        log.Error("failed to get player", sl.Err(err))
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (p *Player) UpdatePlayer(ctx context.Context, uuid string, dto dto.UpdatePlayer) (models.Player, error) {
    const op = "services.player.UpdatePlayer"

    log := p.log.With(slog.String("op", op))
    log.Info("updating player")

    res, err := p.playerUpdater.UpdatePlayer(ctx, uuid, dto)
    if err != nil {
        log.Error("failed to update player", sl.Err(err))
        return models.Player{}, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}