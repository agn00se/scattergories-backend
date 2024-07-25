package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"

	"time"

	"gorm.io/gorm"
)

func StartGame(roomID uint, userID uint) (*models.Game, *models.GameRoomConfig, []*models.GamePrompt, error) {
	// Verify host
	if err := verifyGameRoomHost(roomID, userID, common.ErrStartGameNotHost); err != nil {
		return nil, nil, nil, err
	}

	// Verify no game at the Ongoing or Voting stage
	if err := verifyNoActiveGameInRoom(roomID); err != nil {
		return nil, nil, nil, err
	}

	// Create a new game with the status set to Ongoing
	game := &models.Game{
		GameRoomID: roomID,
		Status:     models.GameStatusOngoing,
		StartTime:  time.Now(),
	}
	if err := repositories.CreateGame(game); err != nil {
		return nil, nil, nil, err
	}

	// Find all users in the GameRoom and create Player entries for the new game
	users, err := getUsersByGameRoomID(roomID)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := createPlayersInGame(users, game.ID); err != nil {
		return nil, nil, nil, err
	}

	// Load GameRoomConfig
	gameRoomConfig, err := getGameRoomConfigByRoomID(roomID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create and load default game prompts
	if err := createGamePrompts(game.ID, gameRoomConfig.NumberOfPrompts); err != nil {
		return nil, nil, nil, err
	}

	gamePrompts, err := getGamePromptsByGameID(game.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	// Return StartGameReponse
	return game, gameRoomConfig, gamePrompts, nil
}

func EndGame(roomID uint, gameID uint, userID uint) (*models.Game, []*models.Player, error) {
	// Verify host
	if err := verifyGameRoomHost(roomID, userID, common.ErrEndGameNotHost); err != nil {
		return nil, nil, err
	}

	// Find the game, set status to completed, and update the end time
	game, err := getGameByID(gameID)
	if err != nil {
		return nil, nil, err
	}
	game.Status = models.GameStatusCompleted
	game.EndTime = time.Now()
	updateGame(game)

	// Calculate final scores
	players, err := getPlayersByGameID(gameID)
	if err != nil {
		return nil, nil, err
	}

	return game, players, nil
}

func getGameByID(gameID uint) (*models.Game, error) {
	return repositories.GetGameByID(gameID)
}

func updateGame(game *models.Game) error {
	return repositories.UpdateGame(game)
}

func verifyNoActiveGameInRoom(roomID uint) error {
	_, err := repositories.GetGameByRoomIDAndStatus(roomID, string(models.GameStatusOngoing), string(models.GameStatusVoting))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No active games found
		}
		return err
	}
	return common.ErrActiveGameExists
}

func getOngoingGameInRoom(roomID uint) (*models.Game, error) {
	game, err := repositories.GetGameByRoomIDAndStatus(roomID, string(models.GameStatusOngoing))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrNoOngoingGameInRoom // No ongoing games found
		}
		return nil, err
	}
	return game, nil
}
