package main

import (
	"net/http"

	"github.com/erfangho/url-shortener/internal/config"
	"github.com/erfangho/url-shortener/internal/handler"
	"github.com/erfangho/url-shortener/internal/repository"
	"github.com/erfangho/url-shortener/internal/routes"
	"github.com/erfangho/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

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
