package routes

import (
	"user/internal/handlers"
	"user/internal/middleware"

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

// NewRouter dynamically sets up routes
func NewRouter(userController *handlers.UserController) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.Logger())

	// Serve static files
	r.Static("/static", "../../static")

	// Define controllers and their routes
	controllers := map[string]Controller{
		"user": {
			Routes: []Route{
				{"GET", "/users/:id", userController.GetUserByID},
				{"POST", "/users", userController.CreateUser},
				{"GET", "/users", userController.ListUsers},
				{"PUT", "/users/:id", userController.UpdateUser},
				{"DELETE", "/users/:id", userController.DeleteUser},
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
