package responses

import "scattergories-backend/internal/models"

type PromptResponse struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`
}

func toPromptResponse(prompt *models.Prompt) *PromptResponse {
	return &PromptResponse{
		ID:   prompt.ID,
		Text: prompt.Text,
	}
}
