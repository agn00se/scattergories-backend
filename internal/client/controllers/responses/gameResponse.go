package responses

import (
	"scattergories-backend/internal/models"
	"time"
)

type GameResponse struct {
	ID         uint              `json:"id"`
	GameRoomID uint              `json:"room_id"`
	Status     models.GameStatus `json:"host_id"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
}

func ToGameResponse(game models.Game) GameResponse {
	return GameResponse{
		ID:         game.ID,
		GameRoomID: game.GameRoomID,
		Status:     game.Status,
		StartTime:  game.StartTime,
		EndTime:    game.EndTime,
	}
}
