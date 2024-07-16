package responses

type StartGameResponse struct {
	Type       string                 `json:"type"`
	Game       GameResponse           `json:"game"`
	GameConfig GameRoomConfigResponse `json:"game_config"`
	Prompts    []GamePromptResponse   `json:"prompts"`
}
