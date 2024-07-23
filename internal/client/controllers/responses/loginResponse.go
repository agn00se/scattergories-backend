package responses

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func ToLoginResponse(accessToken string, refreshToken string) *LoginResponse {
	if refreshToken == "" {
		return &LoginResponse{
			AccessToken: accessToken,
		}
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
