package models

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	Answer       string     `gorm:"not null;size:30" json:"answer"`
	IsValid      bool       `gorm:"not null;default:false" json:"is_valid"`
	PlayerID     uint       `gorm:"not null;uniqueIndex:idx_gameprompt_player" json:"player_id"`
	GamePromptID uint       `gorm:"not null;uniqueIndex:idx_gameprompt_player" json:"game_prompt_id"`
	Player       Player     `gorm:"foreignKey:PlayerID" json:"player"`
	GamePrompt   GamePrompt `gorm:"foreignKey:GamePromptID" json:"game_prompt"`
}
