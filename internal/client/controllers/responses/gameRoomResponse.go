package responses

import "scattergories-backend/internal/models"

type GameRoomResponse struct {
	ID        uint   `json:"id"`
	RoomCode  string `json:"room_code"`
	HostID    uint   `json:"host_id"`
	HostName  string `json:"host_name"`
	IsPrivate bool   `json:"is_private"`
	Passcode  string `json:"passcode,omitempty"`
}

func ToGameRoomResponse(gameRoom models.GameRoom) GameRoomResponse {
	return GameRoomResponse{
		ID:        gameRoom.ID,
		RoomCode:  gameRoom.RoomCode,
		HostID:    gameRoom.HostID,
		HostName:  gameRoom.Host.Name,
		IsPrivate: gameRoom.IsPrivate,
	}
}
