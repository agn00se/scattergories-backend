package requests

type SubmitAnswerRequest struct {
	Type         string `json:"type"`
	Answer       string `json:"answer"`
	UserID       uint   `json:"user_id" validate:"required"`
	GamePromptID uint   `json:"game_prompt_id" validate:"required"`
}
