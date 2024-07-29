package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameRoomDataService interface {
	LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error)
}

type GameRoomDataServiceImpl struct {
	db            *gorm.DB
	answerService AnswerService
	gameService   GameService
}

func NewGameRoomDataService(db *gorm.DB, answerService AnswerService, gameService GameService) GameRoomDataService {
	return &GameRoomDataServiceImpl{
		db:            db,
		answerService: answerService,
		gameService:   gameService,
	}
}

func (s *GameRoomDataServiceImpl) LoadDataForRoom(roomID uuid.UUID) (*domain.Game, []*domain.Answer, error) {
	var game *domain.Game
	var answers []*domain.Answer

	err := utils.WithTransaction(s.db, func(tx *gorm.DB) error {

		// Get the Ongoing game
		var err error
		game, err = s.gameService.GetOngoingGameInRoom(roomID)
		if err != nil {
			return err
		}

		// Set game status to Voting stage and update endtime
		game.Status = domain.GameStatusVoting
		game.EndTime = time.Now()
		if err := s.gameService.UpdateGame(game); err != nil {
			return err
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
