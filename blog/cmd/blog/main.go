package main

import (
	handlers "go-blog/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
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
