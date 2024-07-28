package domain

import "github.com/google/uuid"

type Answer struct {
	BaseModel
	Answer       string     `gorm:"not null;size:30" json:"answer"`
	IsValid      bool       `gorm:"not null;default:false" json:"is_valid"`
	PlayerID     uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_gameprompt_player" json:"player_id"`
	GamePromptID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_gameprompt_player" json:"game_prompt_id"`
	Player       Player     `gorm:"foreignKey:PlayerID" json:"player"`
	GamePrompt   GamePrompt `gorm:"foreignKey:GamePromptID" json:"game_prompt"`
}
