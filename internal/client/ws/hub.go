package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// clients - A map that keeps track of all connected clients. The key is a pointer to a Client struct,
//
//	and the value is a boolean indicating whether the client is active.
//
// broadcast - A channel for broadcasting messages to all client.
// register - A channel for registering new clients.
// unregister - A channel for registering new clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// A global instance of Hub is created with initialized channels and an empty clients map.
var hub = Hub{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

// Upgrade HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
}

// The Run method listens on the register, unregister, and broadcast channels.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // receive a client from the hub.register channel
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// When a message is received on the broadcast channel, it is sent to all registered clients.
			// If a client's send channel is blocked, the client is unregistered, and its send channel is closed.
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Each call to HandleWebSocket handles a new client connection, allowing multiple clients to connect to the server simultaneously.
// Each client gets its own instance of a Client struct, which is then managed by the hub. This allows multiple clients to
// communicate with the server and each other through the WebSocket connection.
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	hub.register <- client // send the client to the hub.register channel

	// writePump is typically run in a separate goroutine to allow the readPump to handle incoming messages immediately.
	go client.writePump()
	client.readPump()
}
