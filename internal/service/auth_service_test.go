package service

import (
	"testing"

	"github.com/erfangho/url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_GenerateToken(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret-key")
	jwtConfig := &config.JWT{}

	authService := NewAuthService(jwtConfig)
	token, err := authService.GenerateToken(1, "test_user")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret-key")
	jwtConfig := &config.JWT{}

	authService := NewAuthService(jwtConfig)
	token, _ := authService.GenerateToken(1, "test_user")
	tokenClaims, err := authService.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, tokenClaims)
	assert.Equal(t, uint(1), tokenClaims.UserID)
	assert.Equal(t, "test_user", tokenClaims.Username)
}

func TestAuthService_TokenFlow(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret-key")

	jwtConfig := &config.JWT{}
	authService := NewAuthService(jwtConfig)

	tests := []struct {
		name     string
		userID   uint
		username string
	}{
		{
			name:     "generate token for user 1",
			userID:   1,
			username: "user1",
		},
		{
			name:     "generate token for user 2",
			userID:   2,
			username: "user2",
		},
		{
			name:     "generate token for another user",
			userID:   100,
			username: "test_user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := authService.GenerateToken(
				tt.userID,
				tt.username,
			)

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, err := authService.ValidateToken(token)

			assert.NoError(t, err)
			assert.NotNil(t, claims)

			if claims != nil {
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.username, claims.Username)
			}
		})
	}
}
