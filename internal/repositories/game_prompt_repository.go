package repositories

import (
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

type GamePromptRepository interface {
	GetGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error)
	GetGameIDByGamePromptID(gamePromptID uint) (uint, error)
	CreateGamePrompt(gamePrompt *domain.GamePrompt) error
}

type GamePromptRepositoryImpl struct {
	db *gorm.DB
}

func NewGamePromptRepository(db *gorm.DB) GamePromptRepository {
	return &GamePromptRepositoryImpl{db: db}
}

func (r *GamePromptRepositoryImpl) GetGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error) {
	var gamePrompts []*domain.GamePrompt
	if err := r.db.Where("game_id = ?", gameID).Preload("Prompt").Find(&gamePrompts).Error; err != nil {
		return nil, err
	}
	return gamePrompts, nil
}

func (r *GamePromptRepositoryImpl) GetGameIDByGamePromptID(gamePromptID uint) (uint, error) {
	var gamePrompt domain.GamePrompt
	if err := r.db.Where("id = ?", gamePromptID).First(&gamePrompt).Error; err != nil {
		return 0, err
	}
	return gamePrompt.GameID, nil
}

func (r *GamePromptRepositoryImpl) CreateGamePrompt(gamePrompt *domain.GamePrompt) error {
	return r.db.Create(gamePrompt).Error
}
