package services

import (
	"scattergories-backend/internal/domain"
	"time"
)

type GameRoomDataService interface {
	LoadDataForRoom(roomID uint) (*domain.Game, []*domain.Answer, error)
}

type GameRoomDataServiceImpl struct {
	answerService AnswerService
	gameService   GameService
}

func NewGameRoomDataService(answerService AnswerService, gameService GameService) GameRoomDataService {
	return &GameRoomDataServiceImpl{
		answerService: answerService,
		gameService:   gameService,
	}
}

func (s *GameRoomDataServiceImpl) LoadDataForRoom(roomID uint) (*domain.Game, []*domain.Answer, error) {
	// Get the Ongoing game
	game, err := s.gameService.GetOngoingGameInRoom(roomID)
	if err != nil {
		return nil, nil, err
	}

	// Set game status to Voting stage and update endtime
	game.Status = domain.GameStatusVoting
	game.EndTime = time.Now()
	if err := s.gameService.UpdateGame(game); err != nil {
		return nil, nil, err
	}

	// Load answers with related Player and GamePrompt (including Prompt)
	answers, err := s.answerService.GetAnswersByGameID(game.ID)
	if err != nil {
		return nil, nil, err
	}

	return game, answers, nil
}
