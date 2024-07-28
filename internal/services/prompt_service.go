package services

import (
	"fmt"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
)

type PromptService interface {
	GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*domain.Prompt, error)
}

type PromptServiceImpl struct {
	promptRepository repositories.PromptRepository
}

func NewPromptService(promptRepository repositories.PromptRepository) PromptService {
	return &PromptServiceImpl{promptRepository: promptRepository}
}

func (s *PromptServiceImpl) GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*domain.Prompt, error) {
	prompts, err := s.promptRepository.GetRandomPromptsGivenLimit(numberOfPrompts)
	if err != nil {
		return nil, err
	}

	// Check if the number of records found is less than the requested limit
	if len(prompts) < numberOfPrompts {
		return nil, fmt.Errorf("only %d prompts found, but %d were requested", len(prompts), numberOfPrompts)
	}

	return prompts, nil
}
