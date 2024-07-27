package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/domain"
)

func GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*domain.Prompt, error) {
	var prompts []*domain.Prompt
	if err := config.DB.Order("RANDOM()").Limit(numberOfPrompts).Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}
