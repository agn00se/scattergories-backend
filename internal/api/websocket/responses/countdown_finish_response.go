package responses

import "scattergories-backend/internal/domain"

type CountdownFinishResponse struct {
	Type    string            `json:"type"`
	Game    *GameResponse     `json:"game"`
	Answers []*AnswerResponse `json:"answers"`
}

func ToCountdownFinishResponse(game *domain.Game, answers []*domain.Answer) *CountdownFinishResponse {
	return &CountdownFinishResponse{
		Type:    "countdown_finish_response",
		Game:    ToGameResponse(game),
		Answers: ToAnswersResponse(answers),
	}
}
