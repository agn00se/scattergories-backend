package ws

import (
	"encoding/json"
	"log"
)

func sendError(client *Client, message string) {
	log.Println("sendError:", message)
	response := map[string]interface{}{
		"type":  "error",
		"error": message,
	}
	sendResponse(client, response)
}

func sendResponse(client *Client, response interface{}) {
	res, err := json.Marshal(response)
	if err != nil {
		log.Println("sendResponse: Error marshaling response:", err)
		return
	}
	client.send <- res
}
