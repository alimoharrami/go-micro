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
	// User Service
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8080"
	}
	userTarget, err := url.Parse(userServiceURL)
	if err != nil {
		log.Fatalf("Invalid user service URL: %v", err)
	}
	userProxy := httputil.NewSingleHostReverseProxy(userTarget)
	userProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", userTarget.Host)
		req.URL.Scheme = userTarget.Scheme
		req.URL.Host = userTarget.Host
	}

	// Auth Service
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8080"
	}
	authTarget, err := url.Parse(authServiceURL)
	if err != nil {
		log.Fatalf("Invalid auth service URL: %v", err)
	}
	authProxy := httputil.NewSingleHostReverseProxy(authTarget)
	authProxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", authTarget.Host)
		req.URL.Scheme = authTarget.Scheme
		req.URL.Host = authTarget.Host
	}

	// 4. Routing
	// Auth routes
	r.Any("/api/auth/*path", func(c *gin.Context) {
		authProxy.ServeHTTP(c.Writer, c.Request)
	})

	// User service catch-all (for now)
	r.Any("/api/*path", func(c *gin.Context) {
		// Avoid double proxying if /api/auth matched (though gin should handle specific first)
		// But strictly speaking, *path catches everything.
		// However, in Gin, specific paths usually take precedence if defined correctly.
		// Let's rely on Gin's matching or checks.
		// Actually, if we define /api/auth/*path, it matches.
		// /api/*path also matches.
		// We can just proxy to user service for everything else.
		userProxy.ServeHTTP(c.Writer, c.Request)
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
