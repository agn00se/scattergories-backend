package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
)

func GetGamePromptsByGameID(gameID uint) ([]*models.GamePrompt, error) {
	return repositories.GetGamePromptsByGameID(gameID)
}

func CreateGamePrompts(gameID uint, numberOfPrompts int) error {
	// Randomly select a subset of prompts
	prompts, err := GetRandomPromptsGivenLimit(numberOfPrompts)
	if err != nil {
		return err
	}

	// Create GamePrompt entries for the selected prompts
	for _, prompt := range prompts {
		gamePrompt := &models.GamePrompt{
			GameID:   gameID,
			PromptID: prompt.ID,
		}
		repositories.CreateGamePrompt(gamePrompt)
	}

	return nil
}
