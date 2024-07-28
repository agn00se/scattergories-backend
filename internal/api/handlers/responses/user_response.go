package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type UserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      *string `json:"email,omitempty"`
	GameRoomID *string `json:"room_id,omitempty"`
}

func ToUserResponse(user *domain.User) *UserResponse {
	var email *string
	if user.Email != nil {
		email = user.Email
	}
	var gameRoomID *string
	if user.GameRoomID != nil {
		str := utils.UUIDToString(*user.GameRoomID)
		gameRoomID = &str
	}
	return &UserResponse{
		ID:         utils.UUIDToString(user.ID),
		Name:       user.Name,
		Email:      email,
		GameRoomID: gameRoomID,
	}
}
