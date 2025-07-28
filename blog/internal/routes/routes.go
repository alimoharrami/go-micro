package routes

import (
	"go-blog/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Route defines the structure for dynamic routing
type Route struct {
	Method     string
	Path       string
	HandleFunc gin.HandlerFunc
}

// Controller defines the structure for a controller with routes
type Controller struct {
	Routes []Route
}

// SetupRouter dynamically sets up routes
func SetRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.LoadHTMLGlob("../../web/*")

	// Serve static files
	r.Static("/static", "../../static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Welcome to Go Backend",
		})
	})

	return r
}
