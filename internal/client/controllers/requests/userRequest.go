package requests

type UserRequest struct {
	Type       string  `json:"type" binding:"required,not_blank,is_valid_user_type"`
	Name       *string `json:"name" binding:"omitempty,name_required_if_registered"`
	Email      *string `json:"email" binding:"omitempty,email,email_required_if_registered"`
	Password   *string `json:"password" binding:"omitempty,password_required_if_registered"`
	GameRoomID *uint   `json:"room_id,omitempty" binding:"omitempty"`
}
