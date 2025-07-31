package routes

import (
	"go-blog/internal/handlers"
	"go-blog/internal/middleware"
	"go-blog/internal/repository"
	"go-blog/internal/service"
	"go-blog/pkg/auth"
	"net/http"

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
func SetRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	r.LoadHTMLGlob("../../web/*")

	// Serve static files
	r.Static("/static", "../../static")

	r.GET("/", auth.JWTAuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Welcome to Go Backend",
		})
	})

	//initialize Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	userService := service.NewUserService(userRepo)
	PostService := service.NewPostService(postRepo)
	AuthService := service.NewAuthService(userRepo)

	userController := handlers.NewUserController(userService)
	postController := handlers.NewPostController(PostService)
	AuthController := handlers.NewAuthController(AuthService)

	r.POST("/api/posts", auth.JWTAuthMiddleware(), postController.CreatePost)

	// Define controllers and their routes
	controllers := map[string]Controller{
		"user": {
			Routes: []Route{
				{"GET", "/users/:id", userController.GetUserByID},
				{"POST", "/users", userController.CreateUser},
			},
		},
		"post": {
			Routes: []Route{
				{"GET", "/posts", postController.ListPosts},
				{"PUT", "/posts/:id", postController.UpdatePost},
				{"DELETE", "/posts/:id", postController.DeletePost},
				{"GET", "/posts/:id", postController.GetPostByID},
			},
		},
		"auth": {
			Routes: []Route{
				{"POST", "/auth", AuthController.Login},
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
