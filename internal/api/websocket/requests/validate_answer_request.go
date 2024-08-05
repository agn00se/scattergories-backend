package requests

type ValidateAnswerRequest struct {
	Type string `json:"type" validate:"required"`
}
