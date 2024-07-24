package responses

import "scattergories-backend/internal/models"

type EndGameResponse struct {
	Type    string            `json:"type"`
	Game    *GameResponse     `json:"game"`
	Players []*PlayerResponse `json:"players"`
}

func ToEndGameResponse(game *models.Game, players []*models.Player) *EndGameResponse {
	return &EndGameResponse{
		Type:    "end_game_response",
		Game:    ToGameResponse(game),
		Players: ToPlayersResponse(players),
	}
}
