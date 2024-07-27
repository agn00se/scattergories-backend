package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/domain"
)

func GetGamePromptsByGameID(gameID uint) ([]*domain.GamePrompt, error) {
	var gamePrompts []*domain.GamePrompt
	if err := config.DB.Where("game_id = ?", gameID).Preload("Prompt").Find(&gamePrompts).Error; err != nil {
		return nil, err
	}
	return gamePrompts, nil
}

func GetGameIDByGamePromptID(gamePromptID uint) (uint, error) {
	var gamePrompt domain.GamePrompt
	if err := config.DB.Where("id = ?", gamePromptID).First(&gamePrompt).Error; err != nil {
		return 0, err
	}
	return gamePrompt.GameID, nil
}

func CreateGamePrompt(gamePrompt *domain.GamePrompt) error {
	return config.DB.Create(gamePrompt).Error
}
