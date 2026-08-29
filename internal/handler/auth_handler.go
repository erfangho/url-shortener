package handler

import (
	"errors"
	"net/http"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "validation error",
				"error":   err.Error(),
			})
			return
		}
		return
	}

	user, err := h.userService.FindByUserName(c.Request.Context(), req.Username)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid credentials",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	err = h.userService.Authenticate(user, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid credentials",
		})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user":    user,
		"token":   token,
	})
}
