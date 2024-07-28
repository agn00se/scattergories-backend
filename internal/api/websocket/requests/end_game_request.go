package requests

type EndGameRequest struct {
	Type   string `json:"type" validate:"required"`
	GameID uint   `json:"game_id" validate:"required"`
}
