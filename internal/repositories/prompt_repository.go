package repositories

import (
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

type PromptRepository interface {
	GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*domain.Prompt, error)
}

type PromptRepositoryImpl struct {
	db *gorm.DB
}

func NewPromptRepository(db *gorm.DB) PromptRepository {
	return &PromptRepositoryImpl{db: db}
}

func (r *PromptRepositoryImpl) GetRandomPromptsGivenLimit(numberOfPrompts int) ([]*domain.Prompt, error) {
	var prompts []*domain.Prompt
	if err := r.db.Order("RANDOM()").Limit(numberOfPrompts).Find(&prompts).Error; err != nil {
		return nil, err
	}
	return prompts, nil
}
