package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"
	"strings"
)

var (
	defaultTimeLimit       = 100
	defaultNumberOfPrompts = 10
)

func UpdateGameConfig(request *models.GameRoomConfig, userID uint) (*models.GameRoomConfig, error) {
	// Verify host
	verifyGameRoomHost(request.GameRoomID, userID, common.ErrUpdateConfigNotHost)

	// Fetch game room config
	gameRoomConfig, err := getGameRoomConfigByRoomID(request.GameRoomID)
	if err != nil {
		return nil, err
	}

	// Update and save game room config
	gameRoomConfig.TimeLimit = request.TimeLimit
	gameRoomConfig.NumberOfPrompts = request.NumberOfPrompts
	gameRoomConfig.Letter = strings.ToUpper(request.Letter)

	if err := repositories.UpdateGameRoomConfig(gameRoomConfig); err != nil {
		return nil, err
	}

	return gameRoomConfig, nil
}

func getGameRoomConfigByRoomID(roomID uint) (*models.GameRoomConfig, error) {
	return repositories.GetGameRoomConfigByRoomID(roomID)
}

func createDefaultGameRoomConfig(gameRoomID uint) error {
	gameRoomConfig := &models.GameRoomConfig{
		GameRoomID:      gameRoomID,
		TimeLimit:       defaultTimeLimit,
		NumberOfPrompts: defaultNumberOfPrompts,
		Letter:          utils.GetRandomLetter(),
	}

	if err := repositories.CreateGameRoomConfig(gameRoomConfig); err != nil {
		return err
	}
	return nil
}
