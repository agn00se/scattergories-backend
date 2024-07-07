package services

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrGameRoomNotFound       = errors.New("game room not found")
	ErrActiveGameExists       = errors.New("active game exists in the room")
	ErrUserNotInSpecifiedRoom = errors.New("user is not in the specified game room")
)
