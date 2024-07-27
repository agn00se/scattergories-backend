package websocket

import (
	"encoding/json"
	"scattergories-backend/internal/api/websocket/requests"
	"scattergories-backend/internal/api/websocket/responses"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/services"
	"scattergories-backend/pkg/validators"
	"time"
)

func HandleMessage(client *Client, roomID uint, messageType string, message []byte) {
	switch messageType {
	case "start_game_request":
		startGame(client, roomID, message)
	case "end_game_request":
		endGame(client, roomID, message)
	case "submit_answer_request":
		submitAnswer(client, roomID, message)
	case "update_game_config_request":
		updateGameConfig(client, roomID, message)
	default:
		sendError(client, "Unknown message type")
	}
}

func startGame(client *Client, roomID uint, message []byte) {
	var req requests.StartGameRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid start_game_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	permitted, err := services.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	// Start the game
	game, gameRoomConfig, gamePrompts, err := services.StartGame(roomID)
	if err != nil {
		sendError(client, err.Error())
		return
	}
	response := responses.ToStartGameResponse(game, gameRoomConfig, gamePrompts)
	sendResponse(client, response)

	// Start the countdown
	countdownDuration := time.Duration(response.GameConfig.TimeLimit) * time.Second
	client.startCountdown(countdownDuration, roomID)
}

func endGame(client *Client, roomID uint, message []byte) {
	var req requests.EndGameRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid end_game_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	permitted, err := services.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	game, players, err := services.EndGame(roomID, req.GameID, req.HostID)
	if err != nil {
		sendError(client, err.Error())
		return
	}
	response := responses.ToEndGameResponse(game, players)
	sendResponse(client, response)
}

func submitAnswer(client *Client, roomID uint, message []byte) {
	var req requests.SubmitAnswerRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid submit_answer_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	if err := services.CreateOrUpdateAnswer(roomID, req.Answer, client.userID, req.GamePromptID); err != nil {
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
		sendError(client, "Invalid update_game_config_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	permitted, err := services.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	config := &domain.GameRoomConfig{
		GameRoomID:      roomID,
		TimeLimit:       req.TimeLimit,
		NumberOfPrompts: req.NumberOfPrompts,
		Letter:          req.Letter,
	}

	config, err = services.UpdateGameConfig(config)
	if err != nil {
		sendError(client, err.Error())
		return
	}

	response := responses.ToUpdateGameConfigResponse(config)
	sendResponse(client, response)
}
