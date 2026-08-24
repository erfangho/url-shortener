package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/erfangho/url-shortener/internal/config"
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/erfangho/url-shortener/internal/repository"
	"github.com/erfangho/url-shortener/internal/routes"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/erfangho/url-shortener/pkg/cache"
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
	logFile, err := os.OpenFile(
		"app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	logger := slog.New(
		slog.NewTextHandler(logFile, nil),
	)

	slog.SetDefault(logger)

	r := gin.Default()
	db, err := config.InitDB()
	urlCache := cache.NewCache(5 * time.Minute)

	if err != nil {
		panic(err)
	}

	urlRepo := repository.NewURLRepository(db, urlCache)
	urlService := service.NewURLService(urlRepo)
	urlHandler := handler.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	routes.RegisterRoutes(r, urlHandler, userHandler)

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
