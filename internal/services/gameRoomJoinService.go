package services

import (
	"scattergories-backend/internal/common"
)

func JoinGameRoom(userID uint, roomID uint) error {
	// Verify game room exists
	_, err := GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify no game at the Ongoing or Voting stage
	if err := VerifyNoActiveGameInRoom(roomID); err != nil {
		return err
	}

	// Verify user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// todo: user limitation - 6 people max

	// Update the associated game room in the user table
	user.GameRoomID = &roomID
	return UpdateUser(user)
}

func LeaveGameRoom(userID uint, roomID uint) error {
	// Verify game room exists
	_, err := GetGameRoomByID(roomID)
	if err != nil {
		return err
	}

	// Verify user exists
	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify user is in the specified game room and remove
	if user.GameRoomID != nil && *user.GameRoomID == roomID {
		user.GameRoomID = nil
	} else {
		return common.ErrUserNotInSpecifiedRoom
	}

	// todo: If host leaves room, assign a new host randomly

	// todo: If last user leaves room, delete game room

	// Update the user record
	return UpdateUser(user)
}
