package services

import (
	"scattergories-backend/internal/common"
)

var roomCapacity = 6

func JoinGameRoom(userID uint, roomID uint) error {
	// Verify game room exists
	_, err := GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify no game at the Ongoing or Voting stage
	if err := verifyNoActiveGameInRoom(roomID); err != nil {
		return err
	}

	// Verify user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify game room not full
	usersInRoom, err := getUsersByGameRoomID(roomID)
	if err != nil {
		return err
	}
	if len(usersInRoom) >= roomCapacity {
		return common.ErrGameRoomFull
	}

	// Update the associated game room in the user table
	user.GameRoomID = &roomID
	return updateUser(user)
}

func LeaveGameRoom(userID uint, roomID uint) error {
	// Verify game room exists
	gameRoom, err := GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify user is in the specified game room and remove game room association
	if user.GameRoomID != nil && *user.GameRoomID == roomID {
		user.GameRoomID = nil
	} else {
		return common.ErrUserNotInSpecifiedRoom
	}

	// Update user record
	if err := updateUser(user); err != nil {
		return err
	}

	// Check if there are any users left in the game room
	usersInRoom, err := getUsersByGameRoomID(roomID)
	if err != nil {
		return err
	}

	// If no users left, delete the game room
	if len(usersInRoom) == 0 {
		if err := DeleteGameRoomByID(roomID); err != nil {
			return err
		}
		return nil
	}

	// Check if the user is the host and assign a new host if needed
	if gameRoom.HostID != nil && *gameRoom.HostID == userID {
		if _, err := updateHost(roomID, usersInRoom[0].ID); err != nil {
			return err
		}
	}

	return nil
}
