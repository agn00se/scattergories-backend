package services

import (
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

func CreateOrUpdateAnswer(roomID uint, answerText string, userID uint, gamePromptID uint) error {
	// Get gameID from gamePromptID
	gameID, err := getGameIDByGamePromptID(gamePromptID)
	if err != nil {
		return err
	}

	// Get the player from userID and gameID
	player, err := getPlayerByUserIDAndGameID(userID, gameID)
	if err != nil {
		return err
	}

	existingAnswer, err := repositories.GetAnswerByPlayerAndPrompt(player.ID, gamePromptID)
	if err == nil {
		// Update the existing answer if one is found
		existingAnswer.Answer = answerText
		return repositories.SaveAnswer(existingAnswer)
	} else if err != gorm.ErrRecordNotFound {
		// Create a new answer if no existing answer is found
		answer := &models.Answer{
			PlayerID:     player.ID,
			GamePromptID: gamePromptID,
			Answer:       answerText,
		}
		return repositories.CreateAnswer(answer)
	} else {
		return err
	}
}

func getAnswersByGameID(gameID uint) ([]*models.Answer, error) {
	return repositories.GetAnswersByGameID(gameID)
}
