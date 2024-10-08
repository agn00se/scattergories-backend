package services

import (
	"scattergories-backend/internal/common"

	"github.com/google/uuid"
)

type PermissionType string

const (
	GameRoomReadPermission  PermissionType = "gameroom:read"
	GameRoomWritePermission PermissionType = "gameroom:write"
	UserWritePermission     PermissionType = "user:write"
)

type PermissionService interface {
	HasPermission(userID uuid.UUID, permissionType PermissionType, resourceID uuid.UUID) (bool, error)
}

type PermissionServiceImpl struct {
	UserService     UserService
	GameRoomService GameRoomService
}

func NewPermissionService(userService UserService, gameRoomService GameRoomService) PermissionService {
	return &PermissionServiceImpl{
		UserService:     userService,
		GameRoomService: gameRoomService,
	}
}

func (s *PermissionServiceImpl) HasPermission(userID uuid.UUID, permissionType PermissionType, resourceID uuid.UUID) (bool, error) {
	switch permissionType {
	case GameRoomReadPermission:
		return s.hasGameRoomReadPermission(userID, resourceID)
	case GameRoomWritePermission:
		return s.hasGameRoomWritePermission(userID, resourceID)
	case UserWritePermission:
		return s.hasUserWritePermission(userID, resourceID)
	}
	return false, nil
}

// GetGameRoom
func (s *PermissionServiceImpl) hasGameRoomReadPermission(userID uuid.UUID, resourceID uuid.UUID) (bool, error) {
	// Verify game room exists
	if _, err := s.GameRoomService.GetGameRoomByID(resourceID); err != nil {
		return false, err
	}

	// Verify user exists
	user, err := s.UserService.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Verify user is in the specified game room
	if user.GameRoomID == nil || *user.GameRoomID != resourceID {
		return false, common.ErrUserNotInSpecifiedRoom
	}
	return true, nil
}

// DeleteGameRoom, StartGame, EndGame, UpdateGameConfig, ValidateAnswers
func (s *PermissionServiceImpl) hasGameRoomWritePermission(userID uuid.UUID, resourceID uuid.UUID) (bool, error) {
	// Verify user exists
	_, err := s.UserService.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	// Verify user is the game room host
	room, err := s.GameRoomService.GetGameRoomByID(resourceID)
	if err != nil {
		return false, err
	}
	if room.HostID != userID {
		return false, common.ErrGameRoomNotHost
	}
	return true, nil
}

// UpdateUser, DeleteUser
func (s *PermissionServiceImpl) hasUserWritePermission(userID uuid.UUID, resourceID uuid.UUID) (bool, error) {
	// Verify user is the user being deleted
	if userID != resourceID {
		return false, common.ErrDeleteUserNotSelf
	}

	return true, nil
}
