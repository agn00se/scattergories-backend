package ws

import (
	"time"

	"github.com/gorilla/websocket"
)

// conn - A pointer to the WebSocket connection.
// send - A channel for sending messages to the client.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Continuously reads messages from the WebSocket connection.
// When a message is received, it is sent to the hub.broadcast channel.
// If an error occurs or the connection is closed, the client is unregistered,
// and the WebSocket connection is closed.
func (c *Client) readPump() {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		hub.broadcast <- message
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			c.conn.WriteMessage(websocket.PingMessage, nil)
		}
	}
}
