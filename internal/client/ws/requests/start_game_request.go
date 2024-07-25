package requests

type StartGameRequest struct {
	Type   string `json:"type" validate:"required"`
	HostID uint   `json:"host_id" validate:"required"`
}
