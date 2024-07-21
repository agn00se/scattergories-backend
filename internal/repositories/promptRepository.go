package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

func GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*models.Prompt, error) {
	var prompts []*models.Prompt
	if err := config.DB.Order("RANDOM()").Limit(numberOfPrompts).Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}
