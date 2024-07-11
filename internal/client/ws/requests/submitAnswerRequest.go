package requests

import "time"

type SubmitAnswerRequest struct {
	Type      string    `json:"type"`
	GameID    uint      `json:"game_id"`
	PromptID  uint      `json:"prompt_id"`
	PlayerID  uint      `json:"player_id"`
	Answer    string    `json:"answer"`
	Timestamp time.Time `json:"timestamp"`
}
