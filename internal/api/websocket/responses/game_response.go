package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
	"time"
)

type GameResponse struct {
	ID         string            `json:"id"`
	GameRoomID string            `json:"room_id"`
	Status     domain.GameStatus `json:"status"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
}

func ToGameResponse(game *domain.Game) *GameResponse {
	return &GameResponse{
		ID:         utils.UUIDToString(game.ID),
		GameRoomID: utils.UUIDToString(game.GameRoomID),
		Status:     game.Status,
		StartTime:  game.StartTime,
		EndTime:    game.EndTime,
	}
}
