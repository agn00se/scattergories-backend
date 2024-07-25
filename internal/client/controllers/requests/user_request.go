package requests

type UserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	GameRoomID *uint  `json:"room_id,omitempty" binding:"omitempty"`
}
