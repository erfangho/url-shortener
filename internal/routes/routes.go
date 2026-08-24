package routes

import (
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, urlHandler *handler.URLHandler) {
	urlGroup := router.Group("/urls")
	{
		urlGroup.GET("", urlHandler.GetAllURLs)
		urlGroup.POST("", urlHandler.CreateURL)
		urlGroup.GET("/:shortCode", urlHandler.GetURL)
	}

	router.GET("/:shortCode", urlHandler.Redirect)
}
