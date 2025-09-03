package routes

import (
	"notification/external/protos/userpb"
	"notification/internal/handlers"
	"notification/internal/middleware"
	"notification/internal/repository"
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
func SetRouter(db *gorm.DB, client userpb.UserServiceClient) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Static("/static", "../../static")

	channelRepo := repository.NewChannelRepository(db)

	channelService := service.NewChannelService(channelRepo, client)

	notificationHand := handlers.NewNotificationController(channelService)
	// Serve static files

	controllers := map[string]Controller{
		"notification": {
			Routes: []Route{
				{"POST", "/notification", notificationHand.Create},
				{"POST", "/notification/subscribe", notificationHand.Subscribe},
				{"POST", "/notification/send", notificationHand.SendNotif},
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
