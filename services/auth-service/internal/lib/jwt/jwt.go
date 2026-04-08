package jwtutil

import (
	"time"

	"github.com/baracudara/hoops/auth-service/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewTokens(
	user models.User, 
	jwtSecret string, 
	accessTokenTTL time.Duration, 
	refreshTokenTTL time.Duration,
	) (accessToken string, refreshToken string, err error) {
		accessToken, err = NewAccessToken(user, jwtSecret, accessTokenTTL)
		if err != nil {
			return "", "", err
		}
		refreshToken, err = NewRefreshToken(user, jwtSecret, refreshTokenTTL)
		if err != nil {
			return "", "", err
		}

		return accessToken, refreshToken, nil
}


func NewAccessToken(user models.User, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["uuid"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(duration).Unix()

	accessToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return accessToken, nil
}


func NewRefreshToken(user models.User, jwtSecret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

    claims["uuid"] = user.ID  
    claims["exp"] = time.Now().Add(duration).Unix()

	refreshToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

