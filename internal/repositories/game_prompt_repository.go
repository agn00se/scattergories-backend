package repositories

import (
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GamePromptRepository interface {
	GetGamePromptsByGameID(gameID uuid.UUID) ([]*domain.GamePrompt, error)
	GetGameIDByGamePromptID(gamePromptID uuid.UUID) (uuid.UUID, error)
	CreateGamePrompt(gamePrompt *domain.GamePrompt) error
}

type GamePromptRepositoryImpl struct {
	db *gorm.DB
}

func NewGamePromptRepository(db *gorm.DB) GamePromptRepository {
	return &GamePromptRepositoryImpl{db: db}
}

func (r *GamePromptRepositoryImpl) GetGamePromptsByGameID(gameID uuid.UUID) ([]*domain.GamePrompt, error) {
	var gamePrompts []*domain.GamePrompt
	if err := r.db.Where("game_id = ?", gameID).Preload("Prompt").Find(&gamePrompts).Error; err != nil {
		return nil, err
	}
	return gamePrompts, nil
}

func (r *GamePromptRepositoryImpl) GetGameIDByGamePromptID(gamePromptID uuid.UUID) (uuid.UUID, error) {
	var gamePrompt domain.GamePrompt
	if err := r.db.Where("id = ?", gamePromptID).First(&gamePrompt).Error; err != nil {
		return uuid.Nil, err
	}
	return gamePrompt.GameID, nil
}

func (r *GamePromptRepositoryImpl) CreateGamePrompt(gamePrompt *domain.GamePrompt) error {
	return r.db.Create(gamePrompt).Error
}
