package websocket

import (
	"log"
	"net/http"
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// Hub maintains the set of active clients and broadcasts messages to the clients in a specific room.
type Hub struct {
	// rooms maps room IDs to a set of active clients in that room. Active client is being kept track by a map of all connected clients
	// where key is a pointer to a Client struct and value is a boolean means active.
	rooms map[uint]map[*Client]bool
	// broadcast is a channel for broadcasting messages to all client in a room
	broadcast chan Message
	// register is a channel for registering new clients.
	register chan *Client
	// unregister is a channel for registering new clients.
	unregister chan *Client
}

// A global instance of Hub is created with initialized channels and an empty clients map.
var HubInstance = &Hub{
	rooms:      make(map[uint]map[*Client]bool),
	broadcast:  make(chan Message),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// Message represents a message to be broadcasted to a room.
type Message struct {
	RoomID  uint
	Content []byte
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
			if _, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
		case client := <-h.unregister:
			if clients, ok := h.rooms[client.roomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
		case message := <-h.broadcast:
			// When a message is received on the broadcast channel, it is sent to all registered clients in the room.
			// If a client's send channel is blocked, the client is unregistered, and its send channel is closed.
			if clients, ok := h.rooms[message.RoomID]; ok {
				for client := range clients {
					select {
					case client.send <- message.Content:
					default:
						close(client.send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.rooms, message.RoomID)
						}
					}
				}
			}
		}
	}
}

// Each call to HandleWebSocket handles a new client connection, allowing multiple clients to connect to the server simultaneously.
// Each client gets its own instance of a Client struct, which is then managed by the hub. This allows multiple clients to
// communicate with the server and each other through the WebSocket connection.
func HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		handlers.HandleError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	roomID, err := handlers.GetIDParam(c, "room_id")
	if err != nil {
		handlers.HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	// Check if the room exists in the database
	_, err = services.GetGameRoomByID(roomID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.HandleError(c, http.StatusNotFound, "Room not found")
		} else {
			handlers.HandleError(c, http.StatusInternalServerError, "Failed to get game room")
		}
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		handlers.HandleError(c, http.StatusInternalServerError, "Failed to upgrade to WebSocket")
		return
	}

	client := &Client{conn: conn, roomID: roomID, userID: userID.(uint), send: make(chan []byte, 256)}
	HubInstance.register <- client // send the client to the hub.register channel

	// writePump is typically run in a separate goroutine to allow the readPump to handle incoming messages immediately.
	go client.writePump()
	go client.readPump()
}
