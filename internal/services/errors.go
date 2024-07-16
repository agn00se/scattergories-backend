package services

import "errors"

var (
	ErrUserNotFound                   = errors.New("user not found")
	ErrGameRoomNotFound               = errors.New("game room not found")
	ErrActiveGameExists               = errors.New("active game exists in the room")
	ErrUserNotInSpecifiedRoom         = errors.New("user is not in the specified game room")
	ErrHostNotFound                   = errors.New("host user not found")
	ErrUserIsAlreadyHostOfAnotherRoom = errors.New("user is already host of another room")
	ErrNotHost                        = errors.New("only host can start a game")
)
