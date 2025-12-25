package routes

import (
	"go-blog/internal/handlers"
	"go-blog/internal/middleware"

	"github.com/alimoharrami/go-micro/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func NewRouter(db *gorm.DB, postController *handlers.PostController) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.Logger())

	r.LoadHTMLGlob("web/*")

	// Serve static files
	r.Static("/static", "../../static")

	r.POST("/api/posts", auth.AuthMiddleware(), auth.RequirePermission("post:create"), postController.CreatePost)

	// Define controllers and their routes
	controllers := map[string]Controller{
		"post": {
			Routes: []Route{
				{"GET", "/posts/:id", postController.GetPostByID},
				{"GET", "/posts", postController.ListPosts},
				{"PUT", "/posts/:id", postController.UpdatePost},
				{"DELETE", "/posts/:id", postController.DeletePost},
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
