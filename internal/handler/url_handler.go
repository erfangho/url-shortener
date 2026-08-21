package handler

import (
	"errors"
	"net/http"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// CreateURL godoc
// @Summary Create a short URL
// @Description Creates a shortened URL from a long URL
// @Tags urls
// @Accept json
// @Produce json
// @Param request body model.CreateURLRequest true "URL to shorten"
// @Success 201 {object} model.URL
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /urls [post]
func (h *URLHandler) CreateURL(c *gin.Context) {
	var req model.CreateURLRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation error",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.service.CreateURL(req.URL)

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

// GetURL godoc
// @Summary Get URL info
// @Description Get information about a shortened URL
// @Tags urls
// @Produce json
// @Param shortCode path string true "Short code"
// @Success 200 {object} model.URL
// @Failure 404 {object} map[string]interface{}
// @Router /urls/{shortCode} [get]
func (h *URLHandler) GetURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "short code is empty",
		})
		return
	}

	result, err := h.service.GetURL(shortCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects to the original URL and tracks the click
// @Tags urls
// @Param shortCode path string true "Short code"
// @Success 301
// @Failure 404 {object} map[string]interface{}
// @Router /{shortCode} [get]
func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "short code is empty",
		})
		return
	}

	result, err := h.service.Redirect(shortCode)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
			"error":   err.Error(),
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, result)
}
