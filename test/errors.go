package test

var (
	ErrActiveGameExists               = "active game exists in the room"
	ErrGameRoomNotFound               = "game room not found"
	ErrUserIsAlreadyHostOfAnotherRoom = "user is already host of another room"
	ErrStartGameNotHost               = "only host can start a game"
	ErrUpdateConfigNotHost            = "only host can update game config"
	ErrUserNotFound                   = "user not found"
	ErrUserNotInSpecifiedRoom         = "user is not in the specified game room"
	ErrInvalidRoomID                  = "Invalid room ID"
	ErrInvalidUserID                  = "Invalid user ID"
	ErrDeleteUserNotSelf              = "user can only delete themselves"
)
