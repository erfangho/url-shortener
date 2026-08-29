package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/erfangho/url-shortener/internal/model"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type URLHandler struct {
	service          *service.URLService
	analyticsService *service.AnalyticsService
}

func NewURLHandler(service *service.URLService, analyticsService *service.AnalyticsService) *URLHandler {
	return &URLHandler{
		service:          service,
		analyticsService: analyticsService,
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
// @Router /u/{shortCode} [get]
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

	clickEvent := model.ClickEvent{
		URLID:     result.ID,
		UserAgent: c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
	}

	h.analyticsService.Publish(clickEvent)

	c.Redirect(http.StatusMovedPermanently, result.OriginalURL)
}

// GetAllURLs godoc
// @Summary Get all URLs
// @Description Get a paginated list of all shortened URLs
// @Tags urls
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Success 200 {object} model.GetAllURLsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /urls [get]
func (h *URLHandler) GetAllURLs(c *gin.Context) {
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

	urls, total, totalPages, err := h.service.GetAllURLs(page, perPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "server error",
		})
		return
	}

	response := model.GetAllURLsResponse{
		URLs: urls,
		Pagination: model.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}
