package repositories

import (
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	GetAnswersByGameID(gameID uuid.UUID) ([]*domain.Answer, error)
	GetAnswerByPlayerAndPrompt(playerID uuid.UUID, gamePromptID uuid.UUID) (*domain.Answer, error)
	SaveAnswer(answer *domain.Answer) error
	CreateAnswer(answer *domain.Answer) error
}

type AnswerRepositoryImpl struct {
	db *gorm.DB
}

func NewAnswerRepository(db *gorm.DB) AnswerRepository {
	return &AnswerRepositoryImpl{db: db}
}

func (r *AnswerRepositoryImpl) GetAnswersByGameID(gameID uuid.UUID) ([]*domain.Answer, error) {
	var answers []*domain.Answer
	if err := r.db.Preload("Player.User").Preload("GamePrompt.Prompt").
		Joins("JOIN players ON players.id = answers.player_id").
		Where("game_prompt_id IN (?)", r.db.Table("game_prompts").Select("id").Where("game_id = ?", gameID)).
		Order("game_prompt_id, players.user_id").
		Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *AnswerRepositoryImpl) GetAnswerByPlayerAndPrompt(playerID uuid.UUID, gamePromptID uuid.UUID) (*domain.Answer, error) {
	var answer domain.Answer
	err := r.db.Where("player_id = ? AND game_prompt_id = ?", playerID, gamePromptID).First(&answer).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (r *AnswerRepositoryImpl) SaveAnswer(answer *domain.Answer) error {
	return r.db.Save(answer).Error
}

func (r *AnswerRepositoryImpl) CreateAnswer(answer *domain.Answer) error {
	return r.db.Create(answer).Error
}
