package requests

type SubmitAnswerRequest struct {
	Type         string `json:"type"`
	Answer       string `json:"answer"`
	PlayerID     uint   `json:"player_id" validate:"required"`
	GamePromptID uint   `json:"game_prompt_id" validate:"required"`
}
