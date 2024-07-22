package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

func GetAnswersByGameID(gameID uint) ([]*models.Answer, error) {
	return repositories.GetAnswersByGameID(gameID)
}

func CreateOrUpdateAnswer(roomID uint, answer *models.Answer) error {
	// Verify player submitting the answer is in the game room
	if err := VerifyPlayerInGameRoom(roomID, answer.PlayerID); err != nil {
		return err
	}

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

func VerifyPlayerInGameRoom(roomID uint, playerID uint) error {
	player, err := repositories.GetPlayerByID(playerID)
	if err != nil {
		return err
	}

	if player.User.GameRoomID == nil || *player.User.GameRoomID != roomID {
		return common.ErrUserNotInSpecifiedRoom
	}

	return nil
}
