package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/erfangho/url-shortener/internal/config"
	grpcclient "github.com/erfangho/url-shortener/internal/grpc"
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/erfangho/url-shortener/internal/middleware"
	"github.com/erfangho/url-shortener/internal/repository"
	"github.com/erfangho/url-shortener/internal/routes"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/erfangho/url-shortener/pkg/cache"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	godotenv.Load()

	redisClient, err := config.NewRedisClient()
	if err != nil {
		panic(err)
	}

	analyticsClient, err := grpcclient.NewAnalyticsClient("localhost:50051")

	if err != nil {
		slog.Error("failed to connect to analytics server", "error", err)
	}

	defer analyticsClient.Close()

	r := gin.Default()
	db, err := config.InitDB()
	urlCache := cache.NewCache(5 * time.Minute)

	if err != nil {
		panic(err)
	}

	urlRepo := repository.NewURLRepository(db, urlCache)
	urlService := service.NewURLService(urlRepo)
	analyticService := service.NewAnalyticsService(urlRepo, 100, 3)
	urlHandler := handler.NewURLHandler(urlService, analyticService)

	jwtConfig := &config.JWT{}
	authService := service.NewAuthService(jwtConfig)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	userRepo := repository.NewUserRepository(db)
	redisCache := cache.NewRedisCache(redisClient, 5*time.Minute)
	userService := service.NewUserService(userRepo, redisCache)
	userHandler := handler.NewUserHandler(userService, authService)

	authHandler := handler.NewAuthHandler(authService, userService)

	routes.RegisterRoutes(r, urlHandler, userHandler, authHandler, authMiddleware)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go func() {
		if err := r.Run(); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	analyticService.Close()
}
