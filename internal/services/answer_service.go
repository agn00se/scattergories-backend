package services

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerService interface {
	CreateOrUpdateAnswer(roomID uuid.UUID, answerText string, userID uuid.UUID, gamePromptID uuid.UUID) error
	GetAnswersByGameID(gameID uuid.UUID) ([]*domain.Answer, error)
}

type AnswerServiceImpl struct {
	db                *gorm.DB
	answerRepository  repositories.AnswerRepository
	gamePromptService GamePromptService
	playerService     PlayerService
}

func NewAnswerService(db *gorm.DB, answerRepository repositories.AnswerRepository, gamePromptService GamePromptService, playerService PlayerService) AnswerService {
	return &AnswerServiceImpl{
		db:                db,
		answerRepository:  answerRepository,
		gamePromptService: gamePromptService,
		playerService:     playerService}
}

func (s *AnswerServiceImpl) CreateOrUpdateAnswer(roomID uuid.UUID, answerText string, userID uuid.UUID, gamePromptID uuid.UUID) error {
	return utils.WithTransaction(s.db, func(tx *gorm.DB) error {

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
	})
}

func (s *AnswerServiceImpl) GetAnswersByGameID(gameID uuid.UUID) ([]*domain.Answer, error) {
	return s.answerRepository.GetAnswersByGameID(gameID)
}
