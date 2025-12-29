package routes

import (
	"notification/external/protos/userpb"
	"notification/internal/handlers"
	"notification/internal/middleware"

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
func SetRouter(db *gorm.DB,
	client userpb.UserServiceClient,
	notificationHandler *handlers.NotificationController,
	notificationBroadCastHandler *handlers.NotificationBroadcastHandler,

) *gin.Engine {
	gin.SetMode("release")

	r := gin.Default()
	r.Use(middleware.Logger())
	r.Static("/static", "../../static")

	// Serve static files
	r.GET("/notification/ws", notificationBroadCastHandler.WebSocket)
	r.POST("/notification/broadcast", notificationBroadCastHandler.HandlePostNotification)

	controllers := map[string]Controller{
		"notification": {
			Routes: []Route{
				{"POST", "/notification", notificationHandler.Create},
				{"POST", "/notification/subscribe", notificationHandler.Subscribe},
				{"POST", "/notification/send", notificationHandler.SendNotif},
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
