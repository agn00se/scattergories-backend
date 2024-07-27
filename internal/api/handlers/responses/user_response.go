package responses

import "scattergories-backend/internal/domain"

type UserResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Email      *string `json:"email,omitempty"`
	GameRoomID *uint   `json:"room_id,omitempty"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	var email *string
	if user.Email != nil {
		email = user.Email
	}
	return &UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      email,
		GameRoomID: user.GameRoomID,
	}
}
