package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"

	"gorm.io/gorm"
)

type AnswerService interface {
	CreateOrUpdateAnswer(roomID uint, answerText string, userID uint, gamePromptID uint) error
	GetAnswersByGameID(gameID uint) ([]*domain.Answer, error)
}

type AnswerServiceImpl struct {
	answerRepository  repositories.AnswerRepository
	gamePromptService GamePromptService
	playerService     PlayerService
}

func NewAnswerService(answerRepository repositories.AnswerRepository, gamePromptService GamePromptService, playerService PlayerService) AnswerService {
	return &AnswerServiceImpl{
		answerRepository:  answerRepository,
		gamePromptService: gamePromptService,
		playerService:     playerService}
}

func (s *AnswerServiceImpl) CreateOrUpdateAnswer(roomID uint, answerText string, userID uint, gamePromptID uint) error {
	// Get gameID from gamePromptID
	gameID, err := s.gamePromptService.getGameIDByGamePromptID(gamePromptID)
	if err != nil {
		return err
	}

	// Get the player from userID and gameID
	player, err := s.playerService.GetPlayerByUserIDAndGameID(userID, gameID)
	if err != nil {
		return err
	}

	existingAnswer, err := s.answerRepository.GetAnswerByPlayerAndPrompt(player.ID, gamePromptID)
	if err == nil {
		// Update the existing answer if one is found
		existingAnswer.Answer = answerText
		return s.answerRepository.SaveAnswer(existingAnswer)
	} else if err == gorm.ErrRecordNotFound {
		// Create a new answer if no existing answer is found
		answer := &domain.Answer{
			PlayerID:     player.ID,
			GamePromptID: gamePromptID,
			Answer:       answerText,
		}
		return s.answerRepository.CreateAnswer(answer)
	} else {
		return err
	}
}

func (s *AnswerServiceImpl) GetAnswersByGameID(gameID uint) ([]*domain.Answer, error) {
	return s.answerRepository.GetAnswersByGameID(gameID)
}
