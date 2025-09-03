package handlers

import (
	"net/http"
	"notification/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	service *service.ChannelService
}

type RequestData struct {
	Channel string `json:"channel"`
	UserID  int    `json:"userID"`
}

type RequestSendNotifData struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

func NewNotificationController(service *service.ChannelService) *NotificationController {
	return &NotificationController{service}
}

func (nc *NotificationController) Create(c *gin.Context) {
	var input service.CreateChannelInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	channel, err := nc.service.Create(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

func (nc *NotificationController) Subscribe(c *gin.Context) {
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}

	errr := nc.service.Subscribe(c, requestData.UserID, requestData.Channel)
	if errr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscribw"})
		return
	}

	c.JSON(http.StatusAccepted, "done")
}

func (nc *NotificationController) SendNotif(c *gin.Context) {
	var requestData RequestSendNotifData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}

	err := nc.service.SendNotif(c, requestData.Channel, requestData.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusAccepted, "done")
}
