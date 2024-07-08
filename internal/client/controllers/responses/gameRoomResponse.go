package responses

import "scattergories-backend/internal/models"

type GameRoomResponse struct {
	ID        uint    `json:"id"`
	RoomCode  string  `json:"room_code"`
	HostID    *uint   `json:"host_id,omitempty"`
	HostName  *string `json:"host_name,omitempty"`
	IsPrivate bool    `json:"is_private"`
}

// Passcode is excluded from return response for security reasons
func ToGameRoomResponse(gameRoom models.GameRoom) GameRoomResponse {
	var hostName *string
	if gameRoom.Host != nil {
		hostName = &gameRoom.Host.Name
	}
	return GameRoomResponse{
		ID:        gameRoom.ID,
		RoomCode:  gameRoom.RoomCode,
		HostID:    gameRoom.HostID,
		HostName:  hostName,
		IsPrivate: gameRoom.IsPrivate,
	}
}
