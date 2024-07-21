package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

func GetAnswersByGameID(gameID uint) ([]*models.Answer, error) {
	return repositories.GetAnswersByGameID(gameID)
}

func CreateOrUpdateAnswer(answer *models.Answer) error {
	existingAnswer, err := repositories.GetAnswerByPlayerAndPrompt(answer.PlayerID, answer.GamePromptID)
	if err == nil {
		existingAnswer.Answer = answer.Answer
		return repositories.SaveAnswer(existingAnswer)
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	// Create a new answer if no existing answer is found
	return repositories.CreateAnswer(answer)
}
