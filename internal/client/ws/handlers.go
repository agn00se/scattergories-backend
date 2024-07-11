package ws

import (
	"github.com/gorilla/websocket"
)

func HandleMessage(conn *websocket.Conn, messageType string, message []byte) {
	switch messageType {
	case "start_game":
		// var req StartGameRequest
		// if err := json.Unmarshal(message, &req); err != nil {
		// 	conn.WriteJSON(map[string]interface{}{
		// 		"type":  "error",
		// 		"error": "Invalid start_game request format",
		// 	})
		// 	return
		// }
		// prompts, err := services.StartGame(req.GameID)
		// if err != nil {
		// 	conn.WriteJSON(map[string]interface{}{
		// 		"type":  "error",
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
		// conn.WriteJSON(map[string]interface{}{
		// 	"type":    "game_started",
		// 	"prompts": prompts,
		// })
	case "submit_answer":
		// var req SubmitAnswerRequest
		// if err := json.Unmarshal(message, &req); err != nil {
		// 	conn.WriteJSON(map[string]interface{}{
		// 		"type":  "error",
		// 		"error": "Invalid submit_answer request format",
		// 	})
		// 	return
		// }
		// // Handle answer submission
	default:
		conn.WriteJSON(map[string]interface{}{
			"type":  "error",
			"error": "Unknown message type",
		})
	}
}
