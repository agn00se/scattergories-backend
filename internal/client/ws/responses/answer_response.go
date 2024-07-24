package responses

import "scattergories-backend/internal/models"

type AnswerResponse struct {
	Answer     string              `json:"answer"`
	IsValid    bool                `json:"is_valid"`
	Player     *PlayerResponse     `json:"player"`
	GamePrompt *GamePromptResponse `json:"game_prompt"`
}

func ToAnswerResponse(answer *models.Answer) *AnswerResponse {
	return &AnswerResponse{
		Answer:     answer.Answer,
		IsValid:    answer.IsValid,
		Player:     ToPlayerResponse(&answer.Player),
		GamePrompt: ToGamePromptResponse(&answer.GamePrompt),
	}
}

func ToAnswersResponse(answers []*models.Answer) []*AnswerResponse {
	responseAnswers := make([]*AnswerResponse, len(answers))
	for i, answer := range answers {
		responseAnswers[i] = ToAnswerResponse(answer)
	}
	return responseAnswers
}
