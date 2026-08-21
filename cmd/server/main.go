package main

import (
	"net/http"

	"github.com/erfangho/url-shortener/internal/config"
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/erfangho/url-shortener/internal/repository"
	"github.com/erfangho/url-shortener/internal/routes"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/erfangho/url-shortener/docs"
)

// @title URL Shortener API
// @version 1.0
// @description A simple URL shortener API built with Gin and GORM
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()
	db, err := config.InitDB()

	if err != nil {
		panic(err)
	}

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	routes.RegisterRoutes(r, urlHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err = r.Run()
	if err != nil {
		return
	}
}
