package rabbitmq

type RequestMessage struct {
	GameID string `json:"game_id"`
	Prompt string `json:"prompt"`
}

type ResponseMessage struct {
	GameID   string `json:"game_id"`
	Response string `json:"response"`
}
