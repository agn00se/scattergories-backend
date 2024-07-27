package responses

import "scattergories-backend/internal/domain"

type StartGameResponse struct {
	Type       string                `json:"type"`
	Game       *GameResponse         `json:"game"`
	GameConfig *GameConfigResponse   `json:"game_config"`
	Prompts    []*GamePromptResponse `json:"prompts"`
}

func ToStartGameResponse(game *domain.Game, gameRoomConfig *domain.GameRoomConfig, gamePrompts []*domain.GamePrompt) *StartGameResponse {
	return &StartGameResponse{
		Type:       "start_game_response",
		Game:       ToGameResponse(game),
		GameConfig: ToGameConfigResponse(gameRoomConfig),
		Prompts:    ToGamePromptsResponse(gamePrompts),
	}
}
