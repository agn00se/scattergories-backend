package requests

type SubmitAnswerRequest struct {
	Type         string `json:"type"`
	Answer       string `json:"answer"`
	GamePromptID string `json:"game_prompt_id" validate:"required"`
}
