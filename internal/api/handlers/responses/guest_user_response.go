package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type GuestUserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	GameRoomID *string `json:"room_id,omitempty"`
	Token      string  `json:"token"`
}

func ToGuestUserResponse(user *domain.User, token string) *GuestUserResponse {
	var gameRoomID *string
	if user.GameRoomID != nil {
		str := utils.UUIDToString(*user.GameRoomID)
		gameRoomID = &str
	}

	return &GuestUserResponse{
		ID:         utils.UUIDToString(user.ID),
		Name:       user.Name,
		GameRoomID: gameRoomID,
		Token:      token,
	}
}
