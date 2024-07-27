package responses

import "scattergories-backend/internal/domain"

type UpdateGameConfigResponse struct {
	Type       string              `json:"type"`
	GameConfig *GameConfigResponse `json:"game_config"`
}

func ToUpdateGameConfigResponse(gameConfig *domain.GameRoomConfig) *UpdateGameConfigResponse {
	return &UpdateGameConfigResponse{
		Type:       "update_game_config_response",
		GameConfig: ToGameConfigResponse(gameConfig),
	}
}
