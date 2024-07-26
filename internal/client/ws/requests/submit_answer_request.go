package requests

type SubmitAnswerRequest struct {
	Type         string `json:"type"`
	Answer       string `json:"answer"`
	GamePromptID uint   `json:"game_prompt_id" validate:"required"`
}
