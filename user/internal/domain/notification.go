package domain

type NotificationData struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

type Notification struct {
	Type string           `json:"type"`
	Data NotificationData `json:"data"`
}
