package responses

type ExchangeTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ToExchangeTokenResponse(accessToken string, refreshToken string) *ExchangeTokenResponse {
	return &ExchangeTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
