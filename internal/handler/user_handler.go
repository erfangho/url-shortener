package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service     *service.UserService
	authService *service.AuthService
}

func NewUserHandler(service *service.UserService, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		service:     service,
		authService: authService,
	}
}

// CreateUser godoc
// @Summary Register a new user
// @Description Create a new user account and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.CreateUserRequest true "User registration data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
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

	token, err := h.authService.GenerateToken(result.ID, result.Username)

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
		"token":   token,
	})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} model.GetAllUsersResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid query parameter",
		})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid query parameter",
		})
		return
	}

	users, total, totalPages, err := h.service.GetAllUsers(c.Request.Context(), page, perPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	response := model.GetAllUsersResponse{
		Users: users,
		Pagination: model.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}
