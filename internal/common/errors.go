package common

import "errors"

var (
	ErrActiveGameExists               = errors.New("active game exists in the room")
	ErrAuthorizationHeaderNotFound    = errors.New("authorization header not found")
	ErrEndGameNotHost                 = errors.New("only host can end the game")
	ErrEmailAlreadyUsed               = errors.New("email is already in use")
	ErrDeleteUserNotSelf              = errors.New("user can only delete themselves")
	ErrGameNotFound                   = errors.New("game not found")
	ErrGamePromptNotFound             = errors.New("game prompt not found")
	ErrGameRoomFull                   = errors.New("game room is full")
	ErrGameRoomNotHost                = errors.New("you are not the host of the game room")
	ErrGameRoomNotFound               = errors.New("game room not found")
	ErrGameRoomConfigNotFound         = errors.New("game room config not found")
	ErrGameRoomWithGivenHostNotFound  = errors.New("game room with the specied host not found")
	ErrInvalidToken                   = errors.New("invalid token")
	ErrLoginFailed                    = errors.New("invalid login credentials")
	ErrNoAnswersToValidate            = errors.New("no answers have been submitted for validation")
	ErrNoOngoingGameInRoom            = errors.New("no ongoing game in room")
	ErrPlayerNotFound                 = errors.New("player not found")
	ErrUserIsAlreadyHostOfAnotherRoom = errors.New("user is already host of another room")
	ErrUserNotFound                   = errors.New("user not found")
	ErrUserNotInSpecifiedRoom         = errors.New("user is not in the specified game room")
)
