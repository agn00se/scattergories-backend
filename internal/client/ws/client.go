package ws

import (
	"encoding/json"
	"log"
	"scattergories-backend/internal/client/ws/responses"
	"scattergories-backend/internal/services"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connected to a specific room.
// conn - A pointer to the WebSocket connection.
// roomID - The ID of the room the client is connected to.
// send - A channel for sending messages to the client.
type Client struct {
	conn   *websocket.Conn
	roomID uint
	send   chan []byte
}

// Continuously reads messages from the WebSocket connection.
// When a message is received, it is sent to the hub.broadcast channel.
// If an error occurs or the connection is closed, the client is unregistered,
// and the WebSocket connection is closed.
func (c *Client) readPump() {
	defer func() {
		HubInstance.unregister <- c
		c.conn.Close()
		log.Printf("Client unregistered from room %d", c.roomID)
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client in room %d: %v", c.roomID, err)
			break
		}

		// Determine message type
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		messageType, ok := msg["type"].(string)
		if !ok {
			c.conn.WriteJSON(map[string]interface{}{
				"type":  "error",
				"error": "Invalid message format: missing type field",
			})
			continue
		}

		HandleMessage(c, c.roomID, messageType, message)

		HubInstance.broadcast <- Message{RoomID: c.roomID, Content: message}
	}
}

// Continuously sends messages from the send channel to the WebSocket connection.
// If the send channel is closed, it sends a close message to the WebSocket connection.
// Sends a ping message every 54 seconds to keep the connection alive.
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The channel is closed, so send a close message to the WebSocket
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// Send a text message to the WebSocket
			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			// Send a ping message to keep the connection alive
			c.conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

// This logic can be moved to client-side
func (c *Client) startCountdown(duration time.Duration, roomID uint) {
	go func() {
		time.Sleep(duration)
		c.triggerWorkflow(roomID)
	}()
}

func (c *Client) triggerWorkflow(roomID uint) {
	game, answers, err := services.LoadDataForRoom(roomID)
	if err != nil {
		sendError(c, "Error loading data")
		return
	}

	response := responses.ToCountdownFinishResponse(game, answers)
	sendResponse(c, response)
}
