package requests

type GameRoomRequest struct {
	IsPrivate bool   `json:"is_private"`
	Passcode  string `json:"passcode,omitempty" binding:"passcode_required_if_private"`
}
