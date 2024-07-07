package requests

type UserRequest struct {
	Name       string `json:"name" binding:"required,not_blank"`
	GameRoomID *uint  `json:"room_id,omitempty"`
}
