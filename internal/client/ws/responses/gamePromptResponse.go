package responses

import "scattergories-backend/internal/models"

type GamePromptResponse struct {
	GameID   uint           `json:"game_id"`
	PromptID uint           `json:"prompt_id"`
	Prompt   PromptResponse `json:"prompt"`
}

func ToGamePromptResponse(gamePrompt models.GamePrompt) GamePromptResponse {
	return GamePromptResponse{
		GameID:   gamePrompt.GameID,
		PromptID: gamePrompt.PromptID,
		Prompt:   toPromptResponse(gamePrompt.Prompt),
	}
}
