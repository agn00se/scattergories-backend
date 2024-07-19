package responses

type StartGameResponse struct {
	Game       GameResponse         `json:"game"`
	GameConfig GameConfigResponse   `json:"game_config"`
	Prompts    []GamePromptResponse `json:"prompts"`
}
