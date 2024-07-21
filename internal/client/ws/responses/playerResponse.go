package responses

import "scattergories-backend/internal/models"

type PlayerResponse struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

func ToPlayerResponse(player *models.Player) *PlayerResponse {
	return &PlayerResponse{
		UserID: player.UserID,
		Name:   player.User.Name,
		Score:  player.Score,
	}
}
