package main

import (
	handlers "go-blog/internal/api"
	"log"

	"github.com/gin-gonic/gin"
	"go-blog/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg.Server.Port)
	r := gin.Default()

	r.LoadHTMLGlob("../../web/*")

	// Routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/post", handlers.ViewPostHandler)
	r.GET("/new", handlers.NewPostHandler)
	r.POST("/create", handlers.CreatePostHandler)

	// Start server
	r.Run(":8080")
}
