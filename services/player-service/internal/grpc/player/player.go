package playergrpc

import (
	"context"

	"errors"

	"github.com/baracudara/hoops/player-service/internal/domain/dto"
	"github.com/baracudara/hoops/player-service/internal/domain/models"
	"github.com/baracudara/hoops/player-service/internal/storage"
	"github.com/baracudara/hoops/protos/gen/go/player"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Player interface {
    CreatePlayer(ctx context.Context, player models.Player) (models.Player, error)
    GetPlayer(ctx context.Context, uuid string) (models.Player, error)
    UpdatePlayer(ctx context.Context, uuid string, dto dto.UpdatePlayer) (models.Player, error)
}

type serverAPI struct {
    player.UnimplementedPlayerServer
    playerService Player
}

func Register(gRPCServer *grpc.Server, playerService Player) {
    player.RegisterPlayerServer(gRPCServer, &serverAPI{playerService: playerService})
}

func (s *serverAPI) GetProfile(ctx context.Context, in *player.GetProfileRequest) (*player.GetProfileResponse, error) {
    if in.GetUuid() == "" {
        return nil, status.Error(codes.InvalidArgument, "uuid is required")
    }

    res, err := s.playerService.GetPlayer(ctx, in.GetUuid())
    if err != nil {
        if errors.Is(err, storage.ErrPlayerNotFound) {
            return nil, status.Error(codes.NotFound, "player not found")
        }
        return nil, status.Error(codes.Internal, "failed to get player")
    }

    return &player.GetProfileResponse{
        Uuid:     res.ID,
        Name:     res.Name,
        Nickname: res.Nickname,
        Position: stringVal(res.Position),
        Age:      int32Val(res.Age),
    }, nil
}

func (s *serverAPI) UpdateProfile(ctx context.Context, in *player.UpdateProfileRequest) (*player.UpdateProfileResponse, error) {
    if in.GetUuid() == "" {
        return nil, status.Error(codes.InvalidArgument, "uuid is required")
    }

    dto := dto.UpdatePlayer{
        Name:     in.GetName(),
        Nickname: in.GetNickname(),
        Position: strPtr(in.GetPosition()),
        Age:      int32Ptr(in.GetAge()),
    }

    _, err := s.playerService.UpdatePlayer(ctx, in.GetUuid(), dto)
    if err != nil {
        if errors.Is(err, storage.ErrPlayerNotFound) {
            return nil, status.Error(codes.NotFound, "player not found")
        }
        return nil, status.Error(codes.Internal, "failed to update player")
    }

    return &player.UpdateProfileResponse{Success: true}, nil
}

func (s *serverAPI) GetStats(ctx context.Context, in *player.GetStatsRequest) (*player.GetStatsResponse, error) {
    if in.GetUuid() == "" {
        return nil, status.Error(codes.InvalidArgument, "uuid is required")
    }

    res, err := s.playerService.GetPlayer(ctx, in.GetUuid())
    if err != nil {
        if errors.Is(err, storage.ErrPlayerNotFound) {
            return nil, status.Error(codes.NotFound, "player not found")
        }
        return nil, status.Error(codes.Internal, "failed to get stats")
    }

    return &player.GetStatsResponse{
        Uuid:        res.ID,
        GamesPlayed: res.Stats.GamesPlayed,
        Wins:        res.Stats.Wins,
        Losses:      res.Stats.Losses,
    }, nil
}


func (s *serverAPI) CreatePlayer(ctx context.Context, in *player.CreatePlayerRequest) (*player.CreatePlayerResponse, error) {
    if in.GetUuid() == "" {
        return nil, status.Error(codes.InvalidArgument, "uuid is required")
    }
    if in.GetName() == "" {
        return nil, status.Error(codes.InvalidArgument, "name is required")
    }
    if in.GetNickname() == "" {
        return nil, status.Error(codes.InvalidArgument, "nickname is required")
    }

    res, err := s.playerService.CreatePlayer(ctx, models.Player{
        ID:       in.GetUuid(),
        Name:     in.GetName(),
        Nickname: in.GetNickname(),
    })
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to create player")
    }

    return &player.CreatePlayerResponse{
        Uuid:     res.ID,
        Name:     res.Name,
        Nickname: res.Nickname,
    }, nil
}

// хелперы для работы с указателями
func stringVal(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

func int32Val(i *int32) int32 {
    if i == nil {
        return 0
    }
    return *i
}

func strPtr(s string) *string {
    if s == "" {
        return nil
    }
    return &s
}

func int32Ptr(i int32) *int32 {
    if i == 0 {
        return nil
    }
    return &i
}