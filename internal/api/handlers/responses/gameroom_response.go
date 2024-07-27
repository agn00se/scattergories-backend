package responses

import "scattergories-backend/internal/domain"

type GameRoomResponse struct {
	ID        uint   `json:"id"`
	RoomCode  string `json:"room_code"`
	HostID    uint   `json:"host_id"`
	HostName  string `json:"host_name"`
	IsPrivate bool   `json:"is_private"`
}

// Passcode is excluded from return response for security reasons
func ToGameRoomResponse(gameRoom *domain.GameRoom) *GameRoomResponse {
	return &GameRoomResponse{
		ID:        gameRoom.ID,
		RoomCode:  gameRoom.RoomCode,
		HostID:    gameRoom.HostID,
		HostName:  gameRoom.Host.Name,
		IsPrivate: gameRoom.IsPrivate,
	}
}
