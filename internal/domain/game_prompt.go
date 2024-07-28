package domain

import "github.com/google/uuid"

// Ensure uniqueness of the combination of GameID and PromptID
type GamePrompt struct {
	BaseModel
	GameID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_game_prompt" json:"game_id"`
	PromptID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_game_prompt" json:"prompt_id"`
	Prompt   Prompt    `gorm:"foreignKey:PromptID"` // Associated Prompt
	Answers  []Answer  `gorm:"foreignKey:GamePromptID;constraint:OnDelete:CASCADE;" json:"-"`
}
