package responses

import "scattergories-backend/internal/domain"

type EndGameResponse struct {
	Type    string            `json:"type"`
	Game    *GameResponse     `json:"game"`
	Players []*PlayerResponse `json:"players"`
}

func ToEndGameResponse(game *domain.Game, players []*domain.Player) *EndGameResponse {
	return &EndGameResponse{
		Type:    "end_game_response",
		Game:    ToGameResponse(game),
		Players: ToPlayersResponse(players),
	}
}
