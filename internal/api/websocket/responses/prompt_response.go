package responses

import "scattergories-backend/internal/domain"

type PromptResponse struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

func toPromptResponse(prompt *domain.Prompt) *PromptResponse {
	return &PromptResponse{
		ID:   prompt.ID,
		Text: prompt.Text,
	}
}
