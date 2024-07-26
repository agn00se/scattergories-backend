package services

import "scattergories-backend/internal/common"

type PermissionType string

const (
	GameRoomReadPermission  PermissionType = "gameroom:read"
	GameRoomWritePermission PermissionType = "gameroom:write"
	UserWritePermission     PermissionType = "user:write"
)

func HasPermission(userID uint, permissionType PermissionType, resourceID uint) (bool, error) {
	switch permissionType {
	case GameRoomReadPermission:
		return hasGameRoomReadPermission(userID, resourceID)
	case GameRoomWritePermission:
		return hasGameRoomWritePermission(userID, resourceID)
	case UserWritePermission:
		return hasUserWritePermission(userID, resourceID)
	}
	return false, nil
}

// GetGameRoom
func hasGameRoomReadPermission(userID uint, resourceID uint) (bool, error) {
	// Verify user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Verify user is in the specified game room
	if user.GameRoomID == nil || *user.GameRoomID != resourceID {
		return false, common.ErrUserNotInSpecifiedRoom
	}
	return true, nil
}

// DeleteGameRoom, StartGame, EndGame, UpdateGameConfig
func hasGameRoomWritePermission(userID uint, resourceID uint) (bool, error) {
	// Verify user exists
	_, err := GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Verify user is the game room host
	room, err := GetGameRoomByID(resourceID)
	if err != nil {
		return false, err
	}
	if *room.HostID != userID {
		return false, common.ErrGameRoomNotHost
	}
	return true, nil
}

// UpdateUser, DeleteUser
func hasUserWritePermission(userID uint, resourceID uint) (bool, error) {
	// Verify user is the user being deleted
	if userID != resourceID {
		return false, common.ErrDeleteUserNotSelf
	}

	return true, nil
}
