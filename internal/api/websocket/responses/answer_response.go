package responses

import "scattergories-backend/internal/domain"

type AnswerResponse struct {
	UserID string `json:"user_id"`
	Answer string `json:"answer"`
}

func ToAnswerResponse(answer *domain.Answer) *AnswerResponse {
	return &AnswerResponse{
		UserID: answer.Player.UserID.String(),
		Answer: answer.Answer,
	}
}

func ToAnswersResponse(answers []*domain.Answer) []*AnswerResponse {
	responseAnswers := make([]*AnswerResponse, len(answers))
	for i, answer := range answers {
		responseAnswers[i] = ToAnswerResponse(answer)
	}
	return responseAnswers
}
