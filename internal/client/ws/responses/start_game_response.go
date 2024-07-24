package responses

import "scattergories-backend/internal/models"

type StartGameResponse struct {
	Type       string                `json:"type"`
	Game       *GameResponse         `json:"game"`
	GameConfig *GameConfigResponse   `json:"game_config"`
	Prompts    []*GamePromptResponse `json:"prompts"`
}

func ToStartGameResponse(game *models.Game, gameRoomConfig *models.GameRoomConfig, gamePrompts []*models.GamePrompt) *StartGameResponse {
	return &StartGameResponse{
		Type:       "start_game_response",
		Game:       ToGameResponse(game),
		GameConfig: ToGameConfigResponse(gameRoomConfig),
		Prompts:    ToGamePromptsResponse(gamePrompts),
	}
}
