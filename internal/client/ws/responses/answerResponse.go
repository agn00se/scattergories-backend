package responses

type AnswerResponse struct {
	Answer     string             `json:"answer"`
	IsValid    bool               `json:"is_valid"`
	Player     PlayerResponse     `json:"player"`
	GamePrompt GamePromptResponse `json:"game_prompt"`
}
