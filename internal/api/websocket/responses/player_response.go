package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type PlayerResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
}

func ToPlayerResponse(player *domain.Player) *PlayerResponse {
	return &PlayerResponse{
		UserID: utils.UUIDToString(player.UserID),
		Name:   player.User.Name,
		Score:  player.Score,
	}
}

func ToPlayersResponse(players []*domain.Player) []*PlayerResponse {
	responsePlayers := make([]*PlayerResponse, len(players))
	for i, player := range players {
		responsePlayers[i] = ToPlayerResponse(player)
	}
	return responsePlayers
}
