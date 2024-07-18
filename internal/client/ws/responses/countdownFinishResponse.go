package responses

type CountdownFinishResponse struct {
	Game    GameResponse     `json:"game"`
	Answers []AnswerResponse `json:"answers"`
}
