package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"strings" // Added missing import

	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// --- Proxy Configuration ---
	userTarget := getProxyTarget("USER_SERVICE_URL", "http://user-service:8080")
	authTarget := getProxyTarget("AUTH_SERVICE_URL", "http://auth-service:8080")
	blogTarget := getProxyTarget("BLOG_SERVICE_URL", "http://blog-service:8080")
	notificationTarget := getProxyTarget("NOTIFICATION_SERVICE_URL", "http://notification-service:8080")

	userProxy := httputil.NewSingleHostReverseProxy(userTarget)
	authProxy := httputil.NewSingleHostReverseProxy(authTarget)
	blogProxy := httputil.NewSingleHostReverseProxy(blogTarget)
	notificationProxy := httputil.NewSingleHostReverseProxy(notificationTarget)

	// --- Routing Logic ---
	r.Any("/api/*path", func(c *gin.Context) {
		path := c.Param("path") // This will be e.g., "/auth/login" or "/users/1"

		// Ensure we are checking prefixes correctly
		if strings.HasPrefix(path, "/auth") {
			authProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		if strings.HasPrefix(path, "/posts") {
			blogProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		if strings.HasPrefix(path, "/notification") {
			notificationProxy.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Default to user service
		userProxy.ServeHTTP(c.Writer, c.Request)
	})

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

// Helper to keep main clean
func getProxyTarget(envVar, defaultURL string) *url.URL {
	uri := os.Getenv(envVar)
	if uri == "" {
		uri = defaultURL
	}
	target, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("Invalid URL for %s: %v", envVar, err)
	}
	return target
}