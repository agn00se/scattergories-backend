package common

import "errors"

var (
	ErrActiveGameExists               = errors.New("active game exists in the room")
	ErrEndGameNotHost                 = errors.New("only host can end the game")
	ErrEmailAlreadyUsed               = errors.New("email is already in use")
	ErrGameNotFound                   = errors.New("game not found")
	ErrGameRoomFull                   = errors.New("game room is full")
	ErrGameRoomNotFound               = errors.New("game room not found")
	ErrGameRoomConfigNotFound         = errors.New("game room config not found")
	ErrGameRoomWithGivenHostNotFound  = errors.New("game room with the specied host not found")
	ErrLoginFailed                    = errors.New("invalid login credentials")
	ErrNoOngoingGameInRoom            = errors.New("no ongoing game in room")
	ErrPlayerNotFound                 = errors.New("player not found")
	ErrStartGameNotHost               = errors.New("only host can start a game")
	ErrUpdateConfigNotHost            = errors.New("only host can update game config")
	ErrUserIsAlreadyHostOfAnotherRoom = errors.New("user is already host of another room")
	ErrUserNotFound                   = errors.New("user not found")
	ErrUserNotInSpecifiedRoom         = errors.New("user is not in the specified game room")
)
