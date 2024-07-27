package requests

type EndGameRequest struct {
	Type   string `json:"type" validate:"required"`
	HostID uint   `json:"host_id" validate:"required"`
	GameID uint   `json:"game_id" validate:"required"`
}
