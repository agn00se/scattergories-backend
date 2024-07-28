package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type GamePromptResponse struct {
	ID     string          `json:"id"`
	Prompt *PromptResponse `json:"prompt"`
}

func ToGamePromptResponse(gamePrompt *domain.GamePrompt) *GamePromptResponse {
	return &GamePromptResponse{
		ID:     utils.UUIDToString(gamePrompt.ID),
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
