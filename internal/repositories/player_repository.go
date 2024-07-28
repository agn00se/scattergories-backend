package repositories

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	GetPlayersByGameID(gameID uuid.UUID) ([]*domain.Player, error)
	GetPlayerByUserIDGameID(userID uuid.UUID, gameID uuid.UUID) (*domain.Player, error)
	CreatePlayer(gamePlayer *domain.Player) error
}

type PlayerRepositoryImpl struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &PlayerRepositoryImpl{db: db}
}

func (r *PlayerRepositoryImpl) GetPlayersByGameID(gameID uuid.UUID) ([]*domain.Player, error) {
	var players []*domain.Player
	if err := r.db.Where("game_id = ?", gameID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func (r *PlayerRepositoryImpl) GetPlayerByUserIDGameID(userID uuid.UUID, gameID uuid.UUID) (*domain.Player, error) {
	var player domain.Player
	if err := r.db.Preload("User").Where("user_id = ? AND game_id = ?", userID, gameID).First(&player).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func (r *PlayerRepositoryImpl) CreatePlayer(gamePlayer *domain.Player) error {
	return r.db.Create(gamePlayer).Error
}
