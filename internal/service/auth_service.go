package service

import (
	"errors"
	"time"

	"github.com/erfangho/url-shortener/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	jwtSecret *config.JWT
}

func NewAuthService(jwtSecret *config.JWT) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

func (a *AuthService) GenerateToken(userId uint, username string) (string, error) {
	claims := &Claims{
		UserID:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.jwtSecret.TokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(a.jwtSecret.JWTSecret())

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return a.jwtSecret.JWTSecret(), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
