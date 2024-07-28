package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

type GamePromptService interface {
	getGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error)
	getGameIDByGamePromptID(gamePromptID uint) (uint, error)
	createGamePrompts(gameID uint, numberOfPrompts int) error
}

type GamePromptServiceImpl struct {
	gamePromptRepository repositories.GamePromptRepository
	promptService        PromptService
}

func NewGamePromptService(gamePromptRepository repositories.GamePromptRepository, promptService PromptService) GamePromptService {
	return &GamePromptServiceImpl{
		gamePromptRepository: gamePromptRepository,
		promptService:        promptService,
	}
}

func (s *GamePromptServiceImpl) getGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error) {
	return s.gamePromptRepository.GetGamePromptsByGameID(gameID)
}

func (s *GamePromptServiceImpl) getGameIDByGamePromptID(gamePromptID uint) (uint, error) {
	gameID, err := s.gamePromptRepository.GetGameIDByGamePromptID(gamePromptID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, common.ErrGamePromptNotFound
		}
		return 0, err
	}
	return gameID, nil
}

func (s *GamePromptServiceImpl) createGamePrompts(gameID uint, numberOfPrompts int) error {
	// Randomly select a subset of prompts
	prompts, err := s.promptService.GetRandomPromptsGivenLimit(numberOfPrompts)
	if err != nil {
		return err
	}

	// Create GamePrompt entries for the selected prompts
	for _, prompt := range prompts {
		gamePrompt := &domain.GamePrompt{
			GameID:   gameID,
			PromptID: prompt.ID,
		}
		s.gamePromptRepository.CreateGamePrompt(gamePrompt)
	}

	return nil
}
