package requests

type StartGameRequest struct {
	Type string `json:"type" validate:"required"`
}
