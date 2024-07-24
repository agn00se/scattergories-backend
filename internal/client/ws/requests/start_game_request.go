package requests

type StartGameRequest struct {
	Type   string `json:"type" validate:"required"`
	UserID uint   `json:"user_id" validate:"required"`
}
