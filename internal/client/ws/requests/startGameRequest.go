package requests

type StartGameRequest struct {
	Type   string `json:"type"`
	GameID uint   `json:"game_id"`
}
