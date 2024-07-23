package requests

type ExchangeTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
