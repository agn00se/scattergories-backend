package repositories

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GameRoomRepository interface {
	GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error)
	GetGameRoomGivenHost(hostID uuid.UUID) (*domain.GameRoom, error)
	GetAllGameRooms() ([]*domain.GameRoom, error)
	CreateGameRoom(gameRoom *domain.GameRoom) error
	UpdateGameRoom(gameRoom *domain.GameRoom) error
	DeleteGameRoomByID(roomID uuid.UUID) *gorm.DB
}

type GameRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewGameRoomRepository(db *gorm.DB) GameRoomRepository {
	return &GameRoomRepositoryImpl{db: db}
}

func (r *GameRoomRepositoryImpl) GetGameRoomByID(roomID uuid.UUID) (*domain.GameRoom, error) {
	var gameRoom domain.GameRoom

	// Eager Preload - tells GORM to load the associated Host object when querying for the GameRoom.
	if err := r.db.Preload("Host").First(&gameRoom, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomNotFound
		}
		return nil, err
	}
	return &gameRoom, nil
}

func (r *GameRoomRepositoryImpl) GetGameRoomGivenHost(hostID uuid.UUID) (*domain.GameRoom, error) {
	var existingGameRoom domain.GameRoom
	if err := r.db.Where("host_id = ?", hostID).First(&existingGameRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomWithGivenHostNotFound
		}
		return nil, err
	}
	return &existingGameRoom, nil
}

func (r *GameRoomRepositoryImpl) GetAllGameRooms() ([]*domain.GameRoom, error) {
	var rooms []*domain.GameRoom
	if err := r.db.Preload("Host").Find(&rooms).Error; err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (r *GameRoomRepositoryImpl) CreateGameRoom(gameRoom *domain.GameRoom) error {
	return r.db.Create(gameRoom).Error
}

func (r *GameRoomRepositoryImpl) UpdateGameRoom(gameRoom *domain.GameRoom) error {
	return r.db.Save(&gameRoom).Error
}

func (r *GameRoomRepositoryImpl) DeleteGameRoomByID(roomID uuid.UUID) *gorm.DB {
	// Delete game room should also delete on cascade any associated
	// game, player, game prompt, game config, and answer from the database
	result := r.db.Unscoped().Delete(&domain.GameRoom{}, roomID)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrGameRoomNotFound
	}
	return result
}
