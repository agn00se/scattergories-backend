package requests

type GameRoomRequest struct {
	HostID    uint   `json:"host_id" binding:"required"`
	IsPrivate bool   `json:"is_private" binding:"required"`
	Passcode  string `json:"passcode,omitempty" binding:"passcode_required_if_private"`
}

type UpdateHostRequest struct {
	NewHostID uint `json:"new_host_id" binding:"required"`
}
