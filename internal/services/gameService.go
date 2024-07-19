package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/ws/responses"
	"scattergories-backend/internal/models"

	"time"

	"gorm.io/gorm"
)

func CreateGame(roomID uint, userID uint) (*responses.StartGameResponse, error) {
	// Check if user is host of the GameRoom
	var gameRoom models.GameRoom
	if err := config.DB.First(&gameRoom, "id = ?", roomID).Error; err != nil {
		return nil, err
	}
	if gameRoom.HostID != nil && *gameRoom.HostID != userID {
		return nil, ErrStartGameNotHost
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
		Game:       responses.ToGameResponse(game),
		GameConfig: responses.ToGameConfigResponse(gameRoomConfig),
		Prompts:    make([]responses.GamePromptResponse, len(gamePrompts)),
	}

	for i, prompt := range gamePrompts {
		response.Prompts[i] = responses.ToGamePromptResponse(prompt)
	}

	return response, nil
}

// todo
func EndGame() {}

func CreateOrUpdateAnswer(answer models.Answer) error {
	var existingAnswer models.Answer
	if err := config.DB.Where("player_id = ? AND game_prompt_id = ?", answer.PlayerID, answer.GamePromptID).First(&existingAnswer).Error; err == nil {
		existingAnswer.Answer = answer.Answer
		return config.DB.Save(&existingAnswer).Error
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	// Create a new answer if no existing answer is found
	return config.DB.Create(&answer).Error
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

func LoadDataForRoom(roomID uint) (*responses.CountdownFinishResponse, error) {
	// Set game status to Voting stage
	var game models.Game
	if err := config.DB.Where("game_room_id = ? AND status = ?", roomID, models.GameStatusOngoing).First(&game).Error; err != nil {
		return nil, err
	}

	game.Status = models.GameStatusVoting
	game.EndTime = time.Now()
	if err := config.DB.Save(&game).Error; err != nil {
		return nil, err
	}

	// Load answers with related Player and GamePrompt (including Prompt)
	var answers []models.Answer
	if err := config.DB.Preload("Player.User").Preload("GamePrompt.Prompt").Where("game_prompt_id IN (?)",
		config.DB.Table("game_prompts").Select("id").Where("game_id = ?", game.ID)).Find(&answers).Error; err != nil {
		return nil, err
	}

	// Map answers to response objects
	responseAnswers := make([]responses.AnswerResponse, len(answers))
	for i, answer := range answers {
		responseAnswers[i] = responses.AnswerResponse{
			Answer:     answer.Answer,
			IsValid:    answer.IsValid,
			Player:     responses.ToPlayerResponse(answer.Player),
			GamePrompt: responses.ToGamePromptResponse(answer.GamePrompt),
		}
	}

	response := &responses.CountdownFinishResponse{
		Game:    responses.ToGameResponse(game),
		Answers: responseAnswers,
	}
	return response, nil
}
