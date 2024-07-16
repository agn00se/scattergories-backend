package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

func CreateGamePrompts(gameID uint, numberOfPrompts int) error {
	// Randomly select a subset of prompts
	var prompts []models.Prompt
	if err := config.DB.Order("RANDOM()").Limit(numberOfPrompts).Find(&prompts).Error; err != nil {
		return err
	}

	// Create GamePrompt entries for the selected prompts
	for _, prompt := range prompts {
		gamePrompt := models.GamePrompt{
			GameID:   gameID,
			PromptID: prompt.ID,
		}
		if err := config.DB.Create(&gamePrompt).Error; err != nil {
			return err
		}
	}

	return nil
}
