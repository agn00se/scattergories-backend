package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func JoinGameRoom(userID uint, roomID uint) error {
	// Check if the game room exists
	gameRoom := models.GameRoom{}
	if err := config.DB.First(&gameRoom, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrGameRoomNotFound
		}
		return err
	}

	// Ensure there isn't any active games in the game room
	var activeGames []models.Game
	if err := config.DB.Where("game_room_id = ? AND (status = ? OR status = ?)", roomID, models.GameStatusOngoing, models.GameStatusVoting).Find(&activeGames).Error; err != nil {
		return err
	}

	if len(activeGames) > 0 {
		return ErrActiveGameExists
	}

	// Check if the user exists
	user := models.User{}
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return err
	}

	// todo: user limitation - 6 people max

	// Assign the user to the game room
	user.GameRoomID = &roomID
	return config.DB.Save(&user).Error
}

func LeaveGameRoom(userID uint, roomID uint) error {
	// Check if the game room exists
	gameRoom := models.GameRoom{}
	if err := config.DB.First(&gameRoom, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrGameRoomNotFound
		}
		return err
	}

	// Check if the user exists
	user := models.User{}
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return err
	}

	// Check if the user is in the specified game room
	if user.GameRoomID != nil && *user.GameRoomID == roomID {
		user.GameRoomID = nil
	} else {
		return ErrUserNotInSpecifiedRoom
	}

	// todo: If host leaves room, assign a new host randomly

	// todo: If last user leaves room, delete game room

	// Update the user record
	return config.DB.Save(&user).Error
}
