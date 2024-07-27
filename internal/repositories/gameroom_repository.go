package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

func GetGameRoomByID(roomID uint) (*domain.GameRoom, error) {
	var gameRoom domain.GameRoom

	// Eager Preload - tells GORM to load the associated Host object when querying for the GameRoom.
	if err := config.DB.Preload("Host").First(&gameRoom, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomNotFound
		}
		return nil, err
	}
	return &gameRoom, nil
}

func GetGameRoomGivenHost(hostID uint) (*domain.GameRoom, error) {
	var existingGameRoom domain.GameRoom
	if err := config.DB.Where("host_id = ?", hostID).First(&existingGameRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomWithGivenHostNotFound
		}
		return nil, err
	}
	return &existingGameRoom, nil
}

func GetAllGameRooms() ([]*domain.GameRoom, error) {
	var rooms []*domain.GameRoom
	if err := config.DB.Preload("Host").Find(&rooms).Error; err != nil {
		return rooms, err
	}
	return rooms, nil
}

func CreateGameRoom(gameRoom *domain.GameRoom) error {
	return config.DB.Create(gameRoom).Error
}

func UpdateGameRoom(gameRoom *domain.GameRoom) error {
	return config.DB.Save(&gameRoom).Error
}

func DeleteGameRoomByID(roomID uint) *gorm.DB {
	// Delete game room should also delete on cascade any associated
	// game, player, game prompt, game config, and answer from the database
	result := config.DB.Unscoped().Delete(&domain.GameRoom{}, roomID)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrGameRoomNotFound
	}
	return result
}
