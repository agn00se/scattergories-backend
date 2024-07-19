package ws

import (
	"encoding/json"
	"scattergories-backend/internal/client/ws/requests"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/services"
	"scattergories-backend/pkg/validators"
	"time"
)

func HandleMessage(client *Client, roomID uint, messageType string, message []byte) {
	switch messageType {
	case "start_game":
		startGame(client, roomID, message)
	case "submit_answer":
		submitAnswer(client, message)
	case "update_game_config":
		updateGameConfig(client, roomID, message)
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
		sendError(client, err.Error())
		return
	}

	sendResponse(client, response)

	// Start the countdown
	countdownDuration := time.Duration(response.GameConfig.TimeLimit) * time.Second
	client.startCountdown(countdownDuration, roomID)
}

func submitAnswer(client *Client, message []byte) {
	var req requests.SubmitAnswerRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid submit_answer request format")
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	answer := models.Answer{
		Answer:       req.Answer,
		PlayerID:     req.PlayerID,
		GamePromptID: req.GamePromptID,
	}

	if err := services.CreateOrUpdateAnswer(answer); err != nil {
		sendError(client, "Failed to save answer"+err.Error())
		return
	}

	sendResponse(client, map[string]interface{}{
		"type":   "submit_answer_response",
		"status": "Answer submitted",
	})
}

func updateGameConfig(client *Client, roomID uint, message []byte) {
	var req requests.UpdateGameConfigRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid submit_answer request format")
		return
	}

	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	config := models.GameRoomConfig{
		GameRoomID:      roomID,
		TimeLimit:       req.TimeLimit,
		NumberOfPrompts: req.NumberOfPrompts,
		Letter:          req.Letter,
	}

	response, err := services.UpdateGameConfig(config, req.UserID)
	if err != nil {
		sendError(client, err.Error())
		return
	}

	sendResponse(client, response)
}
