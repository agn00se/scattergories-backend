package responses

import "scattergories-backend/internal/domain"

type GamePromptResponse struct {
	ID     uint            `json:"id"`
	Prompt *PromptResponse `json:"prompt"`
}

func ToGamePromptResponse(gamePrompt *domain.GamePrompt) *GamePromptResponse {
	return &GamePromptResponse{
		ID:     gamePrompt.ID,
		Prompt: toPromptResponse(&gamePrompt.Prompt),
	}
}

func ToGamePromptsResponse(gamePrompts []*domain.GamePrompt) []*GamePromptResponse {
	responseGamePrompts := make([]*GamePromptResponse, len(gamePrompts))
	for i, gamePrompt := range gamePrompts {
		responseGamePrompts[i] = ToGamePromptResponse(gamePrompt)
	}
	return responseGamePrompts
}
