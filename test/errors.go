package test

var (
	ErrUserNotFound                   = "user not found"
	ErrGameRoomNotFound               = "game room not found"
	ErrActiveGameExists               = "active game exists in the room"
	ErrUserNotInSpecifiedRoom         = "user is not in the specified game room"
	ErrHostNotFound                   = "host user not found"
	ErrUserIsAlreadyHostOfAnotherRoom = "user is already host of another room"
	ErrStartGameNotHost               = "only host can start a game"
	ErrUpdateConfigNotHost            = "only host can update game config"
	ErrInvalidRoomID                  = "Invalid room ID"
	ErrInvalidUserID                  = "Invalid user ID"
)
