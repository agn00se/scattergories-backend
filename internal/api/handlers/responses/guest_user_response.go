package responses

import "scattergories-backend/internal/domain"

type GuestUserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	GameRoomID *uint  `json:"room_id,omitempty"`
	Token      string `json:"token"`
}

func ToGuestUserResponse(user *domain.User, token string) *GuestUserResponse {
	return &GuestUserResponse{
		ID:         user.ID,
		Name:       user.Name,
		GameRoomID: user.GameRoomID,
		Token:      token,
	}
}
