package repositories

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameRepository interface {
	GetGameByID(id uuid.UUID) (*domain.Game, error)
	GetGameByRoomIDAndStatus(roomID uuid.UUID, statuses ...string) (*domain.Game, error)
	CreateGame(game *domain.Game) error
	UpdateGame(game *domain.Game) error
}

type GameRepositoryImpl struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &GameRepositoryImpl{db: db}
}

func (r *GameRepositoryImpl) GetGameByID(id uuid.UUID) (*domain.Game, error) {
	var game domain.Game
	if err := r.db.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameNotFound
		}
		return nil, err
	}
	return &game, nil
}

func (r *GameRepositoryImpl) GetGameByRoomIDAndStatus(roomID uuid.UUID, statuses ...string) (*domain.Game, error) {
	var game domain.Game
	err := r.db.Where("game_room_id = ? AND status IN ?", roomID, statuses).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *GameRepositoryImpl) CreateGame(game *domain.Game) error {
	return r.db.Create(game).Error

}

func (r *GameRepositoryImpl) UpdateGame(game *domain.Game) error {
	return r.db.Save(game).Error
}
