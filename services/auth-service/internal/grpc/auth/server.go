package authgrpc

import (
	"context"
	"errors"

	"github.com/baracudara/hoops/auth-service/internal/domain/dto"
	"github.com/baracudara/hoops/auth-service/internal/domain/models"
	"github.com/baracudara/hoops/auth-service/internal/storage"
	"github.com/baracudara/hoops/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, dto dto.Login) (string, string, error)
	Register(ctx context.Context, dto dto.Register) (string, string, error)
	Logout(ctx context.Context, refreshToken string)  error
	VerifyRefreshToken(ctx context.Context, refreshToken string) (bool, error)
	VerifyAccessToken(ctx context.Context, accessToken string) (models.User, error)
    Refresh(ctx context.Context, refreshToken string) (string, string, error)
}


type serverAPI struct {
	auth.UnimplementedAuthServer
	authService Auth
}

func Register(gRPCServer *grpc.Server, authService Auth) {
    auth.RegisterAuthServer(gRPCServer, &serverAPI{authService: authService})
}


func (s *serverAPI) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {

	if in.GetAuthMethod() == nil {
		return nil, status.Error(codes.InvalidArgument, "auth method is required")
	}

	var d dto.Login

	switch m := in.GetAuthMethod().(type) {
    case *auth.LoginRequest_Email:
			if m.Email.GetEmail() == "" {
                return nil, status.Error(codes.InvalidArgument, "email is required")
        	}
			if m.Email.GetPassword() == "" {
				return nil, status.Error(codes.InvalidArgument, "password is required")
			}
			email := m.Email.GetEmail()
			d.Email = &email
			d.Password = m.Email.GetPassword()
    case *auth.LoginRequest_Phone:
        if m.Phone.GetPhone() == "" {
            return nil, status.Error(codes.InvalidArgument, "phone is required")
        }
        phone := m.Phone.GetPhone()
        d.Phone = &phone

    case *auth.LoginRequest_GoogleId:
        if m.GoogleId == "" {
            return nil, status.Error(codes.InvalidArgument, "google_id is required")
        }
        d.GoogleID = &m.GoogleId

    default:
        return nil, status.Error(codes.InvalidArgument, "unknown auth method")
    }


    accessToken, refreshToken, err := s.authService.Login(ctx, d)
    if err != nil {
        if errors.Is(err, storage.ErrUserExists) {
            return nil, status.Error(codes.AlreadyExists, "user already exists")
        }
        return nil, status.Error(codes.Internal, "failed to register user")
    }

    return &auth.LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}


func (s *serverAPI) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
    if in.GetName() == "" {
        return nil, status.Error(codes.InvalidArgument, "name is required")
    }
    if in.GetNickname() == "" {
        return nil, status.Error(codes.InvalidArgument, "nickname is required")
    }

    // проверяем oneof
    if in.GetAuthMethod() == nil {
        return nil, status.Error(codes.InvalidArgument, "auth method is required")
    }

    d := dto.Register{
        Name:     in.GetName(),
        Nickname: in.GetNickname(),
    }

    // определяем какой метод пришёл
    switch m := in.GetAuthMethod().(type) {
    case *auth.RegisterRequest_Email:
        if m.Email.GetEmail() == "" {
            return nil, status.Error(codes.InvalidArgument, "email is required")
        }
        if m.Email.GetPassword() == "" {
            return nil, status.Error(codes.InvalidArgument, "password is required")
        }
        email := m.Email.GetEmail()
        d.Email = &email
        d.Password = m.Email.GetPassword()

    case *auth.RegisterRequest_Phone:
        if m.Phone.GetPhone() == "" {
            return nil, status.Error(codes.InvalidArgument, "phone is required")
        }
        phone := m.Phone.GetPhone()
        d.Phone = &phone

    case *auth.RegisterRequest_GoogleId:
        if m.GoogleId == "" {
            return nil, status.Error(codes.InvalidArgument, "google_id is required")
        }
        d.GoogleID = &m.GoogleId

    default:
        return nil, status.Error(codes.InvalidArgument, "unknown auth method")
    }

    accessToken, refreshToken, err := s.authService.Register(ctx, d)
    if err != nil {
        if errors.Is(err, storage.ErrUserNotFound) {
            return nil, status.Error(codes.AlreadyExists, "user already exists")
        }
        return nil, status.Error(codes.Internal, "failed to register user")
    }

    return &auth.RegisterResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

func (s *serverAPI) Logout(ctx context.Context, in *auth.LogoutRequest) (*auth.LogoutResponse, error) {
    if in.GetToken() == "" {
        return nil, status.Error(codes.InvalidArgument, "refresh token is required")
    }

    err := s.authService.Logout(ctx, in.GetToken())
    if err != nil {
        if errors.Is(err, storage.ErrTokenNotFound) {
            return &auth.LogoutResponse{Success: false}, nil  
        }
        return nil, status.Error(codes.Internal, "failed to logout")  
    }
    return &auth.LogoutResponse{Success: true}, nil

}


func (s *serverAPI) VerifyRefreshToken(ctx context.Context, in *auth.VerifyRefreshTokenRequest) (*auth.VerifyRefreshTokenResponse, error) {
    if in.GetRefreshToken() == "" {
        return nil, status.Error(codes.InvalidArgument, "refresh token is required")
    }

    valid, err := s.authService.VerifyRefreshToken(ctx, in.GetRefreshToken())

    if err != nil {
        return nil, status.Error(codes.Internal, "failed to virify token")
    }

    return &auth.VerifyRefreshTokenResponse{
        Valid: valid,
    }, nil
}


func (s *serverAPI) VerifyAccessToken(ctx context.Context, in *auth.VerifyAccessTokenRequest) (*auth.VerifyAccessTokenResponse, error) {
    if in.GetAccessToken() == "" {
        return nil, status.Error(codes.InvalidArgument, "refresh token is required")
    }

    user, err := s.authService.VerifyAccessToken(ctx, in.GetAccessToken())

    if err != nil {
        return nil, status.Error(codes.Internal, "failed to virify token")
    }

    return &auth.VerifyAccessTokenResponse{
        Uuid: user.ID,
        Role: string(user.Role),
        Valid: true,
        
    }, nil

}


func (s *serverAPI) Refresh(ctx context.Context, in *auth.RefreshRequest) (*auth.RefreshResponse, error) {
    if in.GetRefreshToken() == "" {
        return nil, status.Error(codes.InvalidArgument, "refresh token is required")
    }

    accessToken, refreshToken, err := s.authService.Refresh(ctx, in.GetRefreshToken())
    if err != nil {
        return nil, status.Error(codes.Internal, "failed to refresh tokens")
    }

    return &auth.RefreshResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}