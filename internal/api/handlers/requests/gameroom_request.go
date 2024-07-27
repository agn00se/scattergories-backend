package requests

type GameRoomRequest struct {
	HostID    uint   `json:"host_id" binding:"required"`
	IsPrivate bool   `json:"is_private"`
	Passcode  string `json:"passcode,omitempty" binding:"passcode_required_if_private"`
}
