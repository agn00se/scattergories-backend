package responses

import "scattergories-backend/internal/models"

type GamePromptResponse struct {
	ID     uint           `json:"id"`
	Prompt PromptResponse `json:"prompt"`
}

func ToGamePromptResponse(gamePrompt models.GamePrompt) GamePromptResponse {
	return GamePromptResponse{
		ID:     gamePrompt.ID,
		Prompt: toPromptResponse(gamePrompt.Prompt),
	}
}
