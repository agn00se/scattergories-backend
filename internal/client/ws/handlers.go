package ws

import (
	"encoding/json"
	"log"
	"scattergories-backend/internal/client/ws/requests"
	"scattergories-backend/internal/services"
	"scattergories-backend/pkg/validators"
)

func HandleMessage(client *Client, roomID uint, messageType string, message []byte) {
	switch messageType {
	case "start_game":
		startGame(client, roomID, message)
	case "submit_answer":
	default:
		sendError(client, "Unknown message type")
	}
}

func startGame(client *Client, roomID uint, message []byte) {
	var req requests.StartGameRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid start_game request format")
		return
	}

	// Validate the request using Gin's binding and validation
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	// Start the game
	response, err := services.CreateGame(roomID, req.UserID)
	if err != nil {
		log.Println("startGame: CreateGame error:", err)
		sendError(client, err.Error())
		return
	}

	sendResponse(client, response)
}
