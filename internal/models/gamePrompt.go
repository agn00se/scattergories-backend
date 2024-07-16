package models

import "gorm.io/gorm"

// Ensure uniqueness of the combination of GameID and PromptID
type GamePrompt struct {
	gorm.Model
	GameID   uint   `gorm:"not null;uniqueIndex:idx_game_prompt" json:"game_id"`
	PromptID uint   `gorm:"not null;uniqueIndex:idx_game_prompt" json:"prompt_id"`
	Prompt   Prompt `gorm:"foreignKey:PromptID"` // Associated Prompt
}
