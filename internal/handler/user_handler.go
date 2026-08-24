package handler

import (
	"errors"
	"net/http"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest

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

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation error",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.service.CreateUser(
		req.Name,
		req.Username,
		req.Password,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data":    result,
	})
}
