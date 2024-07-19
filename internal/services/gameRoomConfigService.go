package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/ws/responses"
	"scattergories-backend/internal/models"
	"scattergories-backend/pkg/utils"
	"strings"
)

var (
	defaultTimeLimit       = 100
	defaultNumberOfPrompts = 10
)

func CreateDefaultGameRoomConfig(gameRoomID uint) error {
	gameRoomConfig := models.GameRoomConfig{
		GameRoomID:      gameRoomID,
		TimeLimit:       defaultTimeLimit,
		NumberOfPrompts: defaultNumberOfPrompts,
		Letter:          utils.GetRandomLetter(),
	}
	if err := config.DB.Create(&gameRoomConfig).Error; err != nil {
		return err
	}
	return nil
}

func UpdateGameConfig(req models.GameRoomConfig, userID uint) (*responses.GameConfigResponse, error) {
	// Verify host
	var gameRoom models.GameRoom
	if err := config.DB.Preload("Host").Where("id = ?", req.GameRoomID).First(&gameRoom).Error; err != nil {
		return nil, err
	}
	if gameRoom.HostID == nil || *gameRoom.HostID != userID {
		return nil, ErrUpdateConfigNotHost
	}

	// Fetch game room config
	var gameRoomConfig models.GameRoomConfig
	if err := config.DB.Where("game_room_id = ?", req.GameRoomID).First(&gameRoomConfig).Error; err != nil {
		return nil, err
	}

	// Update and save game room config
	gameRoomConfig.TimeLimit = req.TimeLimit
	gameRoomConfig.NumberOfPrompts = req.NumberOfPrompts
	gameRoomConfig.Letter = strings.ToUpper(req.Letter)

	if err := config.DB.Save(&gameRoomConfig).Error; err != nil {
		return nil, err
	}

	response := responses.ToGameConfigResponse(gameRoomConfig)

	return &response, nil
}
