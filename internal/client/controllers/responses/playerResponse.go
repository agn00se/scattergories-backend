package responses

import "scattergories-backend/internal/models"

type PlayerResponse struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	GameID uint `json:"game_id"`
	Score  int  `json:"score"`
}

func ToPlayerResponse(player *models.Player) *PlayerResponse {
	return &PlayerResponse{
		UserID: player.UserID,
		GameID: player.GameID,
		Score:  player.Score,
	}
}
