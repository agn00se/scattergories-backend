package responses

import "scattergories-backend/internal/domain"

type GameConfigResponse struct {
	TimeLimit       int    `json:"time_limit"`
	NumberOfPrompts int    `json:"number_of_prompts"`
	Letter          string `json:"letter"`
}

func ToGameConfigResponse(config *domain.GameRoomConfig) *GameConfigResponse {
	return &GameConfigResponse{
		TimeLimit:       config.TimeLimit,
		NumberOfPrompts: config.NumberOfPrompts,
		Letter:          config.Letter,
	}
}
