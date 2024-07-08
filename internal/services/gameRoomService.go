package services

import (
	"fmt"
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
	"scattergories-backend/pkg/utils"
)

func CreateGameRoom(hostID uint, isPrivate bool, passcode string) (models.GameRoom, error) {
	// Verify that the host user exists
	var host models.User
	if err := config.DB.First(&host, hostID).Error; err != nil {
		return models.GameRoom{}, ErrHostNotFound
	}

	// Verify that the host user is not a host in another game room
	var existingGameRoom models.GameRoom
	if err := config.DB.Where("host_id = ?", hostID).First(&existingGameRoom).Error; err == nil {
		return models.GameRoom{}, ErrUserIsAlreadyHostOfAnotherRoom
	}

	gameRoom := models.GameRoom{
		RoomCode:  utils.GenerateRoomCode(),
		HostID:    &hostID,
		IsPrivate: isPrivate,
		Passcode:  passcode,
	}
	if err := config.DB.Create(&gameRoom).Error; err != nil {
		return gameRoom, err
	}

	// Update the user table with the associated game room id
	host.GameRoomID = &gameRoom.ID
	if err := config.DB.Save(&host).Error; err != nil {
		return gameRoom, err
	}

	// Reload the room with the assoicated host
	if err := config.DB.Preload("Host").First(&gameRoom, gameRoom.ID).Error; err != nil {
		return gameRoom, err
	}
	return gameRoom, nil
}

func GetAllGameRooms() ([]models.GameRoom, error) {
	var rooms []models.GameRoom
	if err := config.DB.Preload("Host").Find(&rooms).Error; err != nil {
		return rooms, err
	}
	return rooms, nil
}

func GetGameRoomByID(roomID uint) (models.GameRoom, error) {
	var gameRoom models.GameRoom

	// Eager Preload - tells GORM to load the associated Host object when querying for the GameRoom.
	if err := config.DB.Preload("Host").First(&gameRoom, roomID).Error; err != nil {
		return gameRoom, err
	}
	return gameRoom, nil
}

func DeleteGameRoomByID(roomID uint) error {
	var gameRoom models.GameRoom
	if err := config.DB.First(&gameRoom, roomID).Error; err != nil {
		return err
	}

	if err := config.DB.Unscoped().Delete(&models.GameRoom{}, roomID).Error; err != nil {
		return err
	}
	return nil
}

func UpdateHost(roomID uint, newHostID uint) (models.GameRoom, error) {
	var gameRoom models.GameRoom
	if err := config.DB.First(&gameRoom, roomID).Error; err != nil {
		return gameRoom, err
	}

	var existingGameRoom models.GameRoom
	if err := config.DB.Where("host_id = ?", newHostID).First(&existingGameRoom).Error; err == nil {
		return models.GameRoom{}, fmt.Errorf("user is already a host in another game room")
	}

	gameRoom.HostID = &newHostID
	if err := config.DB.Save(&gameRoom).Error; err != nil {
		return gameRoom, err
	}

	// Reload the room with the assoicated host
	if err := config.DB.Preload("Host").First(&gameRoom, gameRoom.ID).Error; err != nil {
		return gameRoom, err
	}
	return gameRoom, nil
}
