package responses

import "scattergories-backend/internal/models"

type CountdownFinishResponse struct {
	Type    string            `json:"type"`
	Game    *GameResponse     `json:"game"`
	Answers []*AnswerResponse `json:"answers"`
}

func ToCountdownFinishResponse(game *models.Game, answers []*models.Answer) *CountdownFinishResponse {
	return &CountdownFinishResponse{
		Type:    "countdown_finish_response",
		Game:    ToGameResponse(game),
		Answers: ToAnswersResponse(answers),
	}
}
