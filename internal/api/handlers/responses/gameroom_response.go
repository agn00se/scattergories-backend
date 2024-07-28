package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type GameRoomResponse struct {
	ID        string `json:"id"`
	RoomCode  string `json:"room_code"`
	HostID    string `json:"host_id"`
	HostName  string `json:"host_name"`
	IsPrivate bool   `json:"is_private"`
}

// Passcode is excluded from return response for security reasons
func ToGameRoomResponse(gameRoom *domain.GameRoom) *GameRoomResponse {
	return &GameRoomResponse{
		ID:        utils.UUIDToString(gameRoom.ID),
		RoomCode:  gameRoom.RoomCode,
		HostID:    utils.UUIDToString(gameRoom.HostID),
		HostName:  gameRoom.Host.Name,
		IsPrivate: gameRoom.IsPrivate,
	}
}
