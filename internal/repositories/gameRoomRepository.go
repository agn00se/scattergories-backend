package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func GetGameRoomByID(roomID uint) (*models.GameRoom, error) {
	var gameRoom models.GameRoom

	// Eager Preload - tells GORM to load the associated Host object when querying for the GameRoom.
	if err := config.DB.Preload("Host").First(&gameRoom, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomNotFound
		}
		return nil, err
	}
	return &gameRoom, nil
}

func GetGameRoomGivenHost(hostID uint) (*models.GameRoom, error) {
	var existingGameRoom models.GameRoom
	if err := config.DB.Where("host_id = ?", hostID).First(&existingGameRoom).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomWithGivenHostNotFound
		}
		return nil, err
	}
	return &existingGameRoom, nil
}

func GetAllGameRooms() ([]*models.GameRoom, error) {
	var rooms []*models.GameRoom
	if err := config.DB.Preload("Host").Find(&rooms).Error; err != nil {
		return rooms, err
	}
	return rooms, nil
}

func CreateGameRoom(gameRoom *models.GameRoom) error {
	return config.DB.Create(gameRoom).Error
}

func UpdateGameRoom(gameRoom *models.GameRoom) error {
	return config.DB.Save(&gameRoom).Error
}

func DeleteGameRoomByID(roomID uint) *gorm.DB {
	// Delete game room should also delete on cascade any associated
	// game, player, game prompt, game config, and answer from the database
	result := config.DB.Unscoped().Delete(&models.GameRoom{}, roomID)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrGameRoomNotFound
	}
	return result
}
