package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/ws/responses"
	"scattergories-backend/internal/models"

	"time"
)

func CreateGame(roomID uint, userID uint) (*responses.StartGameResponse, error) {
	// Check if user is host of the GameRoom
	var gameRoom models.GameRoom
	if err := config.DB.First(&gameRoom, "id = ?", roomID).Error; err != nil {
		return nil, err
	}
	if gameRoom.HostID != nil && *gameRoom.HostID != userID {
		return nil, ErrNotHost
	}

	// Check if the GameRoom has any ongoing or voting games
	var existingGames []models.Game
	if err := config.DB.Where("game_room_id = ? AND (status = ? OR status = ?)", roomID, models.GameStatusOngoing, models.GameStatusVoting).Find(&existingGames).Error; err != nil {
		return nil, err
	}

	if len(existingGames) > 0 {
		return nil, ErrActiveGameExists
	}

	// Create the new game with the status set to Ongoing
	game := models.Game{
		GameRoomID: roomID,
		Status:     models.GameStatusOngoing,
		StartTime:  time.Now(),
	}
	if err := config.DB.Create(&game).Error; err != nil {
		return nil, err
	}

	// Find all users in the GameRoom and create Player entries for the new game
	var users []models.User
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&users).Error; err != nil {
		return nil, err
	}

	for _, user := range users {
		gamePlayer := models.Player{
			UserID: user.ID,
			GameID: game.ID,
			Score:  0,
		}
		if err := config.DB.Create(&gamePlayer).Error; err != nil {
			return nil, err
		}
	}

	// Load GameRoomConfig
	var gameRoomConfig models.GameRoomConfig
	if err := config.DB.First(&gameRoomConfig, "game_room_id = ?", roomID).Error; err != nil {
		return nil, err
	}

	// Create default game prompts
	if err := CreateGamePrompts(game.ID, gameRoomConfig.NumberOfPrompts); err != nil {
		return nil, err
	}

	// Load GamePrompt
	var gamePrompts []models.GamePrompt
	if err := config.DB.Where("game_id = ?", game.ID).Preload("Prompt").Find(&gamePrompts).Error; err != nil {
		return nil, err
	}

	response := &responses.StartGameResponse{
		Type:       "start_game_response",
		Game:       responses.ToGameResponse(game),
		GameConfig: responses.ToGameRoomConfigResponse(gameRoomConfig),
		Prompts:    make([]responses.GamePromptResponse, len(gamePrompts)),
	}

	for i, prompt := range gamePrompts {
		response.Prompts[i] = responses.ToGamePromptResponse(prompt)
	}

	return response, nil
}

func GetGamesByRoomID(roomID uint) ([]models.Game, error) {
	var games []models.Game
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&games).Error; err != nil {
		return games, err
	}
	return games, nil
}

func GetGameByID(roomID uint, gameID uint) (models.Game, error) {
	var game models.Game
	if err := config.DB.Where("id = ? AND game_room_id = ?", gameID, roomID).First(&game).Error; err != nil {
		return models.Game{}, err
	}
	return game, nil
}
