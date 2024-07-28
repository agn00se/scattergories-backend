package repositories

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

type GameRoomConfigRepository interface {
	GetGameRoomConfigByRoomID(roomID uint) (*domain.GameRoomConfig, error)
	CreateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error
	UpdateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error
}

type GameRoomConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewGameRoomConfigRepository(db *gorm.DB) GameRoomConfigRepository {
	return &GameRoomConfigRepositoryImpl{db: db}
}

func (r *GameRoomConfigRepositoryImpl) GetGameRoomConfigByRoomID(roomID uint) (*domain.GameRoomConfig, error) {
	var gameRoomConfig domain.GameRoomConfig
	if err := r.db.First(&gameRoomConfig, "game_room_id = ?", roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomConfigNotFound
		}
		return nil, err
	}
	return &gameRoomConfig, nil
}

func (r *GameRoomConfigRepositoryImpl) CreateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error {
	return r.db.Create(gameRoomConfig).Error
}

func (r *GameRoomConfigRepositoryImpl) UpdateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error {
	return r.db.Save(&gameRoomConfig).Error
}
