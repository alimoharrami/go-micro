package hub

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type NotificationHelper struct {
	clients   map[*websocket.Conn]bool
	broadcast chan string
	mu        sync.Mutex
}

func NewNotificationHelper() *NotificationHelper {
	h := &NotificationHelper{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
	}
	go h.run()
	return h
}

func (h *NotificationHelper) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
}

func (h *NotificationHelper) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	conn.Close()
}

func (h *NotificationHelper) Broadcast(msg string) {
	h.broadcast <- msg
}

func (h *NotificationHelper) run() {
	for msg := range h.broadcast {
		var deadClients []*websocket.Conn

		h.mu.Lock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				fmt.Println("Write error:", err)
				deadClients = append(deadClients, client)
			}
		}
		h.mu.Unlock()

		for _, c := range deadClients {
			h.RemoveClient(c)
		}
	}
}