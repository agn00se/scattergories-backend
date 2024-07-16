package responses

import "scattergories-backend/internal/models"

type GameRoomConfigResponse struct {
	GameRoomID      uint   `json:"game_room_id"`
	TimeLimit       int    `json:"time_limit"`
	NumberOfPrompts int    `json:"number_of_prompts"`
	Letter          string `json:"letter"`
}

func ToGameRoomConfigResponse(config models.GameRoomConfig) GameRoomConfigResponse {
	return GameRoomConfigResponse{
		GameRoomID:      config.GameRoomID,
		TimeLimit:       config.TimeLimit,
		NumberOfPrompts: config.NumberOfPrompts,
		Letter:          config.Letter,
	}
}
