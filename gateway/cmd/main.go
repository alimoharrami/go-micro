package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// 3. Proxy Configuration
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8080"
	}

	target, err := url.Parse(userServiceURL)
	if err != nil {
		log.Fatalf("Invalid user service URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	// 4. Routing
	// Forward everything under /api/* to the backend
	r.Any("/api/*path", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	
	// Create a catch-all that returns 404 for non-api routes, 
	// since we are no longer serving static files.
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "Endpoint not found, this is the API Gateway"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
