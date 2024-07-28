package responses

import (
	"scattergories-backend/internal/domain"
	"scattergories-backend/pkg/utils"
)

type PromptResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func toPromptResponse(prompt *domain.Prompt) *PromptResponse {
	return &PromptResponse{
		ID:   utils.UUIDToString(prompt.ID),
		Text: prompt.Text,
	}
}
