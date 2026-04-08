package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/baracudara/hoops/auth-service/internal/domain/dto"
	"github.com/baracudara/hoops/auth-service/internal/domain/models"
	jwtutil "github.com/baracudara/hoops/auth-service/internal/lib/jwt"
	"github.com/baracudara/hoops/auth-service/internal/lib/logger/sl"
	"github.com/baracudara/hoops/auth-service/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials")
)

type Auth struct {
	log *slog.Logger
	usrSaver UserSaver
	usrProvider UserProvider 
	tokenSaver TokenSaver
	tokenChecker TokenChecker
	tokenDeletr TokenDeletr
	accessTokenTTL time.Duration
	refreshTokenTTL time.Duration
	jwtSecret string
}

func New(
	log *slog.Logger,
	usrSaver UserSaver, 
	usrProvider UserProvider, 
	tokenSaver TokenSaver,
	tokenChecker TokenChecker,
	tokenDeletr TokenDeletr,
	accessTokenTTL time.Duration, 
	refreshTokenTTL time.Duration, 
	jwtSecret string,
) *Auth {
	return &Auth{
		log: log, 
		usrSaver: usrSaver,
		usrProvider: usrProvider,
		tokenSaver: tokenSaver,
		tokenChecker: tokenChecker,
		tokenDeletr: tokenDeletr,
		accessTokenTTL: accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		jwtSecret: jwtSecret,

	}

}

type TokenSaver interface {
	SaveToken(
		ctx context.Context, 
		uuid string, 
		token string, 
		ttl time.Duration,
	) error
}

type TokenChecker interface {
	IsTokenValid(
		ctx context.Context, 
		token string,
	) (bool, error)
}

type TokenDeletr interface {
	DeleteToken(
		ctx context.Context, 
		token string,
	) error
}

type UserSaver interface {
	SaveUser(
		ctx context.Context, 
		user models.User, 
	) (models.User, error) 
}

type UserProvider interface {
	GetUser(
		ctx context.Context, 
		dto dto.Login,
	) (models.User, error)
}

func (a *Auth) Login(ctx context.Context, dto dto.Login) (string, string, error) {
	const op = "services.auth.login"

	log := a.log.With(
		slog.String("op", op), 
	)

	log.Info("logging in a user")

	user, err := a.usrProvider.GetUser(ctx, dto)
    if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
        log.Error("failed to get user", sl.Err(err))
        return "", "", fmt.Errorf("%s: %w", op, err)
    }

	if dto.Password != "" {
		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(dto.Password)); err !=  nil {
				log.Warn("invalid password", sl.Err(err))
        		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}

	accessToken, refreshToken, err := jwtutil.NewTokens(user, a.jwtSecret, a.accessTokenTTL, a.refreshTokenTTL)

	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.tokenSaver.SaveToken(ctx, user.ID, refreshToken, a.refreshTokenTTL)

	if err != nil {
		log.Error("failed to save token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken,refreshToken, nil
}

func (a *Auth) Register(ctx context.Context, dto dto.Register) (string, string, error) {
	const op = "services.auth.register"

	log := a.log.With(
		slog.String("op", op), 
	)

	log.Info("registering a user")


	user := models.User{
        ID:          uuid.New().String(), 
		Name:        dto.Name,
		Nickname:    dto.Nickname,
		Email:       dto.Email,
		Phone:       dto.Phone,
		GoogleID:    dto.GoogleID,
		Role:        models.RolePlayer, 
		TrustRating: 100,    
    }

	if dto.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, err)
		}
		user.PassHash = hash
	}

	usersql, err := a.usrSaver.SaveUser(ctx, user)

	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, refreshToken, err := jwtutil.NewTokens(usersql, a.jwtSecret, a.accessTokenTTL, a.refreshTokenTTL)

	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.tokenSaver.SaveToken(ctx, usersql.ID, refreshToken, a.refreshTokenTTL)

	if err != nil {
		log.Error("failed to save token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken,refreshToken, nil
}

func (a *Auth) Logout(ctx context.Context, refreshToken string)  error {
	const op = "services.auth.logout" 

	log := a.log.With(slog.String("op", op))
	log.Info("logging out user")

	err := a.tokenDeletr.DeleteToken(ctx, refreshToken)

	if err != nil {
		log.Error("failed to delete token", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}


func (a *Auth) VerifyRefreshToken(ctx context.Context, refreshToken string) (bool, error) {
	const op = "services.auth.verify.refresh"

	log := a.log.With(
		slog.String("op", op), 
	)

	log.Info("verifiying refresh token")


	res, err := a.tokenChecker.IsTokenValid(ctx, refreshToken)

	if err != nil {
		log.Error("falied to verify token", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil

}

func (a *Auth) VerifyAccessToken(ctx context.Context, accessToken string) (models.User, error) {
	const op = "services.auth.verify.access"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("verifiying access token")

    token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, 
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.jwtSecret), nil
	})
	

	    if err != nil || !token.Valid {
        return models.User{}, fmt.Errorf("%s: invalid token", op)
    }

    claims := token.Claims.(*jwt.MapClaims)
    
    return models.User{
        ID:   (*claims)["uuid"].(string),
        Role: models.Role((*claims)["role"].(string)),
    }, nil
}