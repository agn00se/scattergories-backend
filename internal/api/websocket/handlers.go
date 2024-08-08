package websocket

import (
	"encoding/json"
	"scattergories-backend/internal/api/websocket/requests"
	"scattergories-backend/internal/api/websocket/responses"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/services"
	"scattergories-backend/pkg/utils"
	"scattergories-backend/pkg/validators"
	"time"

	"github.com/google/uuid"
)

type MessageHandler interface {
	HandleMessage(client *Client, roomID uuid.UUID, messageType string, message []byte)
	GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error)
	LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error)
}

type MessageHandlerImpl struct {
	gameService             services.GameService
	gameRoomService         services.GameRoomService
	gameRoomDataService     services.GameRoomDataService
	permissionService       services.PermissionService
	answerService           services.AnswerService
	gameConfigService       services.GameConfigService
	answerValidationService services.AnswerValidationService
}

func NewMessageHandler(
	gameService services.GameService,
	gameRoomService services.GameRoomService,
	gameRoomDataService services.GameRoomDataService,
	permissionService services.PermissionService,
	answerService services.AnswerService,
	gameConfigService services.GameConfigService,
	answerValidationService services.AnswerValidationService,
) MessageHandler {
	return &MessageHandlerImpl{
		gameService:             gameService,
		gameRoomService:         gameRoomService,
		gameRoomDataService:     gameRoomDataService,
		permissionService:       permissionService,
		answerService:           answerService,
		gameConfigService:       gameConfigService,
		answerValidationService: answerValidationService,
	}
}

func (h *MessageHandlerImpl) HandleMessage(client *Client, roomID uuid.UUID, messageType string, message []byte) {
	switch messageType {
	case "start_game_request":
		h.startGame(client, roomID)
	case "end_game_request":
		h.endGame(client, roomID, message)
	case "submit_answer_request":
		h.submitAnswer(client, roomID, message)
	case "update_game_config_request":
		h.updateGameConfig(client, roomID, message)
	case "validate_answers_request":
		h.validateAnswers(client, roomID)
	default:
		sendError(client, "Unknown message type")
	}
}

func (h *MessageHandlerImpl) GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error) {
	return h.gameRoomService.GetGameRoomByID(roomID)
}

func (h *MessageHandlerImpl) LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error) {
	return h.gameRoomDataService.LoadDataForRoom(roomID)
}

func (h *MessageHandlerImpl) startGame(client *Client, roomID uuid.UUID) {
	permitted, err := h.permissionService.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	game, gameRoomConfig, gamePrompts, err := h.gameService.StartGame(roomID)
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

func (h *MessageHandlerImpl) endGame(client *Client, roomID uuid.UUID, message []byte) {
	var req requests.EndGameRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid end_game_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	gameID, err := utils.StringToUUID(req.GameID)
	if err != nil {
		sendError(client, "Invalid game ID format")
		return
	}

	permitted, err := h.permissionService.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	game, players, err := h.gameService.EndGame(roomID, gameID)
	if err != nil {
		sendError(client, err.Error())
		return
	}
	response := responses.ToEndGameResponse(game, players)
	sendResponse(client, response)
}

func (h *MessageHandlerImpl) submitAnswer(client *Client, roomID uuid.UUID, message []byte) {
	var req requests.SubmitAnswerRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid submit_answer_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	gamePromptID, err := utils.StringToUUID(req.GamePromptID)
	if err != nil {
		sendError(client, "Invalid game prompt ID format")
		return
	}

	if err := h.answerService.CreateOrUpdateAnswer(roomID, req.Answer, client.userID, gamePromptID); err != nil {
		sendError(client, "Failed to save answer"+err.Error())
		return
	}

	sendResponse(client, map[string]interface{}{
		"type":   "submit_answer_response",
		"status": "Answer submitted",
	})
}

func (h *MessageHandlerImpl) updateGameConfig(client *Client, roomID uuid.UUID, message []byte) {
	var req requests.UpdateGameConfigRequest
	if err := json.Unmarshal(message, &req); err != nil {
		sendError(client, "Invalid update_game_config_request format")
		return
	}
	if err := validators.Validate.Struct(req); err != nil {
		sendError(client, "Validation failed: "+err.Error())
		return
	}

	permitted, err := h.permissionService.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
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

	config, err = h.gameConfigService.UpdateGameConfig(config)
	if err != nil {
		sendError(client, err.Error())
		return
	}

	response := responses.ToUpdateGameConfigResponse(config)
	sendResponse(client, response)
}

func (h *MessageHandlerImpl) validateAnswers(client *Client, roomID uuid.UUID) {
	permitted, err := h.permissionService.HasPermission(client.userID, services.GameRoomWritePermission, roomID)
	if err != nil || !permitted {
		sendError(client, err.Error())
		return
	}

	if err := h.answerValidationService.ValidateAnswers(roomID); err != nil {
		sendError(client, err.Error())
		return
	}
	sendResponse(client, map[string]interface{}{
		"type":   "validate_answers_response",
		"status": "Answer validated",
	})
}
