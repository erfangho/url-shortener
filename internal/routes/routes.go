package routes

import (
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/erfangho/url-shortener/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, urlHandler *handler.URLHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	urlGroup := router.Group("/urls")
	{
		urlGroup.GET("", authMiddleware.AuthMiddleware(), urlHandler.GetAllURLs)
		urlGroup.POST("", urlHandler.CreateURL)
		urlGroup.GET("/:shortCode", urlHandler.GetURL)
	}

	router.GET("/:shortCode", urlHandler.Redirect)

	userGroup := router.Group("/users")
	{
		userGroup.POST("", userHandler.CreateUser)
	}
}
