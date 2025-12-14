package routes

import (
	"auth/internal/handlers"
	"auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Route defines the structure for dynamic routing
type Route struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

// Controller defines the structure for a controller with routes
type Controller struct {
	Routes []Route
}

// SetupRouter dynamically sets up routes
// NewRouter dynamically sets up routes
func NewRouter(authController *handlers.AuthController) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Serve static files
	r.Static("/static", "../../static")

	// Define controllers and their routes
	controllers := map[string]Controller{
		"auth": {
			Routes: []Route{
				{"POST", "/auth", authController.Login},
			},
		},
	}

	// Register all routes dynamically
	api := r.Group("/api")
	for _, controller := range controllers {
		for _, route := range controller.Routes {
			api.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}

	return r
}
