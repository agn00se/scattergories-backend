package responses

import "scattergories-backend/internal/models"

type GameRoomConfigResponse struct {
	TimeLimit       int    `json:"time_limit"`
	NumberOfPrompts int    `json:"number_of_prompts"`
	Letter          string `json:"letter"`
}

func ToGameRoomConfigResponse(config models.GameRoomConfig) GameRoomConfigResponse {
	return GameRoomConfigResponse{
		TimeLimit:       config.TimeLimit,
		NumberOfPrompts: config.NumberOfPrompts,
		Letter:          config.Letter,
	}
}
