package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/domain"
)

func GetAnswersByGameID(gameID uint) ([]*domain.Answer, error) {
	var answers []*domain.Answer
	if err := config.DB.Preload("Player.User").Preload("GamePrompt.Prompt").Where("game_prompt_id IN (?)",
		config.DB.Table("game_prompts").Select("id").Where("game_id = ?", gameID)).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func GetAnswerByPlayerAndPrompt(playerID uint, gamePromptID uint) (*domain.Answer, error) {
	var answer domain.Answer
	err := config.DB.Where("player_id = ? AND game_prompt_id = ?", playerID, gamePromptID).First(&answer).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func SaveAnswer(answer *domain.Answer) error {
	return config.DB.Save(answer).Error
}

func CreateAnswer(answer *domain.Answer) error {
	return config.DB.Create(answer).Error
}
