package requests

type StartGameResponse struct {
	Type   string `json:"type"`
	GameID uint   `json:"game_id"`
}
