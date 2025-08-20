package routes

import (
	"notification/internal/helpers"
	"notification/internal/middleware"
	"notification/internal/service"

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

	// Serve static files
	r.Static("/static", "../../static")

	//grpc service
	client := helpers.InitGRPC()

	service.NewEmailService(client)

	return r
}
