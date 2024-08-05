package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameRoomDataService interface {
	LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error)
	GetAnswersToBeValidated(roomID uuid.UUID) ([]map[string]interface{}, error)
}

type GameRoomDataServiceImpl struct {
	db                *gorm.DB
	answerService     AnswerService
	gameService       GameService
	playerService     PlayerService
	gamePromptService GamePromptService
}

func NewGameRoomDataService(
	db *gorm.DB,
	answerService AnswerService,
	gameService GameService,
	playerService PlayerService,
	gamePromptService GamePromptService,
) GameRoomDataService {
	return &GameRoomDataServiceImpl{
		db:                db,
		answerService:     answerService,
		gameService:       gameService,
		playerService:     playerService,
		gamePromptService: gamePromptService,
	}
}

func (s *GameRoomDataServiceImpl) LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error) {
	var game *domain.Game
	var answers []*domain.Answer

	err := utils.WithTransaction(s.db, func(tx *gorm.DB) error {

		// Get the Ongoing game
		var err error
		game, err = s.gameService.GetGameByRoomIDAndStatus(roomID, domain.GameStatusOngoing)
		if err != nil {
			return err
		}

		// Set game status to Voting stage and update endtime
		if game.Status == domain.GameStatusOngoing {
			game.Status = domain.GameStatusVoting
			game.EndTime = time.Now()
			if err := s.gameService.UpdateGame(game); err != nil {
				return err
			}
		}

		// Load answers with related Player and GamePrompt (including Prompt)
		answers, err = s.answerService.GetAnswersByGameID(game.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	return game, answers, nil
}

func (s *GameRoomDataServiceImpl) GetAnswersToBeValidated(roomID uuid.UUID) ([]map[string]interface{}, error) {
	game, err := s.gameService.GetGameByRoomIDAndStatus(roomID, domain.GameStatusVoting)
	if err != nil {
		return nil, err
	}

	gamePrompts, err := s.gamePromptService.GetGamePromptsByGameID(game.ID)
	if err != nil {
		return nil, err
	}

	players, err := s.playerService.GetPlayersByGameID(game.ID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}

	for _, gamePrompt := range gamePrompts {
		answers := make([]map[string]interface{}, len(players))
		playerAnswerMap := make(map[uuid.UUID]string)

		// Map answers to player IDs
		for _, ans := range gamePrompt.Answers {
			playerAnswerMap[ans.Player.UserID] = ans.Answer
		}

		// Ensure each player has an entry, even if they didn't submit an answer
		for i, player := range players {
			answer, ok := playerAnswerMap[player.UserID]
			if !ok {
				answer = "" // Placeholder for no answer
			}
			answers[i] = map[string]interface{}{
				"player_id": player.ID,
				"answer":    answer,
			}
		}

		gamePromptResponse := map[string]interface{}{
			"game_prompt_id": gamePrompt.ID,
			"game_prompt":    gamePrompt.Prompt.Text,
			"answers":        answers,
		}
		response = append(response, gamePromptResponse)
	}

	// Handle the case where no answers have been submitted by any player
	if len(response) == 0 {
		return nil, common.ErrNoAnswersToValidate
	}

	return response, nil
}
