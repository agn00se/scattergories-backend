package responses

import "scattergories-backend/internal/models"

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	GameRoomID *uint  `json:"room_id,omitempty"`
}

func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      *user.Email,
		GameRoomID: user.GameRoomID,
	}
}