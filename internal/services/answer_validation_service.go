package services

import (
	"encoding/json"
	"fmt"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/rabbitmq"
	"strings"

	"github.com/google/uuid"
)

type AnswerValidationService interface {
	ValidateAnswers(roomID uuid.UUID) error
	SanitizeAnswers(promptAnswers []map[string]interface{}, letter string) []map[string]interface{}
	ValidateAnswersWithLLM(gameID uuid.UUID, promptAnswers []map[string]interface{}) error
	ConstructPrompt(promptAnswers []map[string]interface{}) string
	ProcessLLMResponse(gameID string, response string) error
	ParseLLMResponse(responseText string, promptAnswers []map[string]interface{}) (map[uuid.UUID][]map[string]interface{}, error)
}

type AnswerValidationServiceImpl struct {
	rabbitMQ            *rabbitmq.RabbitMQ
	gameRoomDataService GameRoomDataService
	gameConfigService   GameConfigService
	gameService         GameService
	answerService       AnswerService
}

func NewAnswerValidationService(
	rabbitMQ *rabbitmq.RabbitMQ,
	gameRoomDataService GameRoomDataService,
	gameConfigService GameConfigService,
	gameService GameService,
	answerService AnswerService,
) AnswerValidationService {
	return &AnswerValidationServiceImpl{
		rabbitMQ:            rabbitMQ,
		gameRoomDataService: gameRoomDataService,
		gameConfigService:   gameConfigService,
		gameService:         gameService,
		answerService:       answerService,
	}
}

const promptInstructions = `Evaluate the following answers to determine if they are valid based on the given context. 
Answer each with 'Yes' or 'No'. If the answer is empty, return 'No'. Provide the responses in the exact format: "Prompt X - Answer Y: Yes/No".

`

// Function to validate answers using OpenAI's GPT-4
func (s *AnswerValidationServiceImpl) ValidateAnswers(roomID uuid.UUID) error {
	// Get the ongoing game in the game room
	game, err := s.gameService.GetGameByRoomIDAndStatus(roomID, domain.GameStatusVoting)
	if err != nil {
		return err
	}

	// Load prompts and promptAnswers for the game
	promptAnswers, err := s.gameRoomDataService.GetAnswersToBeValidated(game.ID)
	if err != nil {
		return err
	}

	// Get game letter
	config, err := s.gameConfigService.GetGameRoomConfigByRoomID(roomID)
	if err != nil {
		return err
	}

	// Sanitize and map answers to a map of prompt:answers
	sanitizedPromptAnswers := s.SanitizeAnswers(promptAnswers, config.Letter)

	// Validate answers using the LLM via message queue
	err = s.ValidateAnswersWithLLM(game.ID, sanitizedPromptAnswers)
	if err != nil {
		return err
	}

	return nil
}

/*

[[ Prompt Answers Format ]]
[
	{
		"game_prompt_id": uuid.New(),
		"game_prompt":    "Name a fruit",
		"answers": []map[string]interface{}{
			{"player_id": uuid.New(), "answer": "Apple"},
			{"player_id": uuid.New(), "answer": "Banana"},
			{"player_id": uuid.New(), "answer": "Cherry"},
		},
	},
	{
		"game_prompt_id": uuid.New(),
		"game_prompt":    "Name a color",
		"answers": []map[string]interface{}{
			{"player_id": uuid.New(), "answer": "Red"},
			{"player_id": uuid.New(), "answer": "Blue"},
			{"player_id": uuid.New(), "answer": "Green"},
		},
	},
]
*/

func (s *AnswerValidationServiceImpl) SanitizeAnswers(promptAnswers []map[string]interface{}, letter string) []map[string]interface{} {
	for _, answer := range promptAnswers {
		answersList := answer["answers"].([]map[string]interface{})
		for _, ans := range answersList {
			answerText := ans["answer"].(string)
			if !strings.HasPrefix(strings.ToLower(answerText), strings.ToLower(letter)) {
				ans["answer"] = "" // Replace invalid answers with an empty string
			}
		}
	}
	return promptAnswers
}

func (s *AnswerValidationServiceImpl) ValidateAnswersWithLLM(gameID uuid.UUID, promptAnswers []map[string]interface{}) error {
	// Construct LLM Prompt
	prompt := s.ConstructPrompt(promptAnswers)

	// Create the request message
	reqMsg := rabbitmq.RequestMessage{
		GameID: gameID.String(),
		Prompt: prompt,
	}

	// Marshal the request message
	messageBody, err := json.Marshal(reqMsg)
	if err != nil {
		return err
	}

	// Send the request message to the message queue
	err = s.rabbitMQ.Publish("llm_request_queue", messageBody)
	if err != nil {
		return err
	}

	return nil
}

/*
[[ LLM Prompt Format ]]

Prompt 1 - Movie Title
1. 'E.T.'
2. 'Endgame'
3. ”
4. ”
5. 'Ex Machina'
6. 'EA Sports'

Prompt 2 - Historical Figure
1. 'Einstein'
2. 'Edison'
3. 'Eminem'
4. ”
5. 'Elizabeth I'
6. ”
*/
func (s *AnswerValidationServiceImpl) ConstructPrompt(promptAnswers []map[string]interface{}) string {
	var promptBuilder strings.Builder
	promptBuilder.WriteString(promptInstructions)

	promptNumber := 1
	for _, promptAnswer := range promptAnswers {
		promptText := promptAnswer["game_prompt"].(string)
		promptBuilder.WriteString(fmt.Sprintf("Prompt %d - %s\n", promptNumber, promptText))
		answersList := promptAnswer["answers"].([]map[string]interface{})
		for i, ans := range answersList {
			answerText := ans["answer"].(string)
			promptBuilder.WriteString(fmt.Sprintf("%d. '%s'\n", i+1, answerText))
		}
		promptBuilder.WriteString("\n")
		promptNumber++
	}
	return promptBuilder.String()
}

func (s *AnswerValidationServiceImpl) ProcessLLMResponse(gameIDString string, response string) error {
	gameID, err := uuid.Parse(gameIDString)
	if err != nil {
		return fmt.Errorf("invalid gameID: %v", err)
	}

	promptAnswers, err := s.gameRoomDataService.GetAnswersToBeValidated(gameID)
	if err != nil {
		return err
	}

	results, err := s.ParseLLMResponse(response, promptAnswers)
	if err != nil {
		return err
	}

	// Update Answer table
	for gamePromptID, validations := range results {
		for _, validation := range validations {
			playerID := validation["player_id"].(uuid.UUID)
			valid := validation["valid"].(bool)

			// No need to update validity for invalid answers since the field is initialized as false
			if valid {
				if err := s.answerService.SetAnswerValid(playerID, gamePromptID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

/*
[[ LLM Response Format ]]

Prompt 1 - Answer 1: Yes
Prompt 1 - Answer 2: Yes
Prompt 1 - Answer 3: No
Prompt 1 - Answer 4: No
Prompt 1 - Answer 5: Yes
Prompt 1 - Answer 6: No

Prompt 2 - Answer 1: Yes
Prompt 2 - Answer 2: Yes
Prompt 2 - Answer 3: No
Prompt 2 - Answer 4: No
Prompt 2 - Answer 5: Yes
Prompt 2 - Answer 6: No

[[ Parsed Result Format ]]

	{
		"prompt-id-1": [
			{"player_id": "player-id-1", "valid": true},
			{"player_id": "player-id-2", "valid": false},
			{"player_id": "player-id-3", "valid": true},
		],
		"prompt-id-2": [
			{"player_id": "player-id-1", "valid": false},
			{"player_id": "player-id-2", "valid": true},
			{"player_id": "player-id-3", "valid": true},
		],
	}
*/
func (s *AnswerValidationServiceImpl) ParseLLMResponse(responseText string, promptAnswers []map[string]interface{}) (map[uuid.UUID][]map[string]interface{}, error) {
	lines := strings.Split(responseText, "\n")
	results := make(map[uuid.UUID][]map[string]interface{})
	var currentPromptID uuid.UUID
	var currentPromptResults []map[string]interface{}
	promptIndex := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if len(currentPromptResults) > 0 {
				results[currentPromptID] = currentPromptResults
				currentPromptResults = nil
			}
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], "Prompt") {
			currentPromptID = promptAnswers[promptIndex]["game_prompt_id"].(uuid.UUID)
			promptIndex++
			continue
		}

		parts = strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format in response line: %s", line)
		}

		answerPart := strings.TrimSpace(parts[1])
		valid := false
		if answerPart == "Yes" {
			valid = true
		} else if answerPart == "No" {
			valid = false
		} else {
			return nil, fmt.Errorf("unexpected answer format in response line: %s", line)
		}

		answersList := promptAnswers[promptIndex-1]["answers"].([]map[string]interface{})
		playerID := answersList[len(currentPromptResults)]["player_id"]

		currentPromptResults = append(currentPromptResults, map[string]interface{}{
			"player_id": playerID,
			"valid":     valid,
		})
	}

	if len(currentPromptResults) > 0 {
		results[currentPromptID] = currentPromptResults
	}

	return results, nil
}
