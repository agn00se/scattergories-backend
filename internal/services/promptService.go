package services

import (
	"fmt"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
)

func GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*models.Prompt, error) {
	prompts, err := repositories.GetRandomPromptsGivenLimit(numberOfPrompts)
	if err != nil {
		return nil, err
	}

	// Check if the number of records found is less than the requested limit
	if len(prompts) < numberOfPrompts {
		return nil, fmt.Errorf("only %d prompts found, but %d were requested", len(prompts), numberOfPrompts)
	}

	return prompts, nil
}
