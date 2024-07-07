package responses

import "scattergories-backend/internal/models"

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	GameRoomID *uint  `json:"room_id,omitempty"`
}

func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		GameRoomID: user.GameRoomID,
	}
}
