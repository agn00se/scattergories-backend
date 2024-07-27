package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

func getGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error) {
	return repositories.GetGamePromptsByGameID(gameID)
}

func getGameIDByGamePromptID(gamePromptID uint) (uint, error) {
	gameID, err := repositories.GetGameIDByGamePromptID(gamePromptID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, common.ErrGamePromptNotFound
		}
		return 0, err
	}
	return gameID, nil
}

func createGamePrompts(gameID uint, numberOfPrompts int) error {
	// Randomly select a subset of prompts
	prompts, err := getRandomPromptsGivenLimit(numberOfPrompts)
	if err != nil {
		return err
	}

	// Create GamePrompt entries for the selected prompts
	for _, prompt := range prompts {
		gamePrompt := &domain.GamePrompt{
			GameID:   gameID,
			PromptID: prompt.ID,
		}
		repositories.CreateGamePrompt(gamePrompt)
	}

	return nil
}
