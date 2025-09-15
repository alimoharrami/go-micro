package handlers

import (
	"fmt"
	"net/http"
	"notification/internal/hub"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type NotificationBroadcastHandler struct {
	hub *hub.NotificationHelper
}

func NewNotificationBoradcastHandler(hub *hub.NotificationHelper) *NotificationBroadcastHandler {
	return &NotificationBroadcastHandler{hub: hub}
}

func (h *NotificationBroadcastHandler) WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	h.hub.AddClient(conn)

	// Listen for client messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.hub.RemoveClient(conn)
			break
		}
		h.hub.Broadcast(string(msg))
	}
}

func (h *NotificationBroadcastHandler) HandlePostNotification(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	// Broadcast the message
	h.hub.Broadcast(req.Message)

	c.JSON(http.StatusOK, gin.H{"status": "sent", "message": req.Message})
}
