package requests

type EndGameRequest struct {
	Type   string `json:"type" validate:"required"`
	GameID string `json:"game_id" validate:"required"`
}
