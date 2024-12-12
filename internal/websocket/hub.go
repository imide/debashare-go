package websocket

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Client represents a WebSocket client connection
type Client struct {
	Conn     *websocket.Conn
	RoomID   string
	Messages chan []byte
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by room
	clients map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations on clients map
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, exists := h.clients[client.RoomID]; !exists {
				h.clients[client.RoomID] = make(map[*Client]bool)
			}
			h.clients[client.RoomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, exists := h.clients[client.RoomID]; exists {
				if _, ok := h.clients[client.RoomID][client]; ok {
					delete(h.clients[client.RoomID], client)
					close(client.Messages)
					if len(h.clients[client.RoomID]) == 0 {
						delete(h.clients, client.RoomID)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Broadcast sends a message to all clients in a room
func (h *Hub) Broadcast(roomID string, message []byte) {
	h.mu.RLock()
	if clients, exists := h.clients[roomID]; exists {
		for client := range clients {
			select {
			case client.Messages <- message:
			default:
				h.mu.RUnlock()
				h.unregister <- client
				h.mu.RLock()
			}
		}
	}
	h.mu.RUnlock()
}

// Register adds a new client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}
