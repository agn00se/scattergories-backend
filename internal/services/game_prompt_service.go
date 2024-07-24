package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
)

func getGamePromptsByGameID(gameID uint) ([]*models.GamePrompt, error) {
	return repositories.GetGamePromptsByGameID(gameID)
}

func createGamePrompts(gameID uint, numberOfPrompts int) error {
	// Randomly select a subset of prompts
	prompts, err := getRandomPromptsGivenLimit(numberOfPrompts)
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
