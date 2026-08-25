package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		parts := strings.SplitN(authorization, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid auth header",
			})
			return
		}

		token := parts[1]

		claims, err := m.authService.ValidateToken(token)

		if err != nil {
			slog.Error("token validate error", "error", err)

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
