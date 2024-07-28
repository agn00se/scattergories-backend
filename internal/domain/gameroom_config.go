package domain

import "github.com/google/uuid"

type GameRoomConfig struct {
	BaseModel
	GameRoomID      uuid.UUID `gorm:"type:uuid;not null;index" json:"game_room_id"`
	TimeLimit       int       `gorm:"not null" json:"time_limit"`          // Time limit for the game in minutes
	NumberOfPrompts int       `gorm:"not null" json:"number_of_prompts"`   // Number of prompts to be used in the game
	Letter          string    `gorm:"type:char(1);not null" json:"letter"` // Store as a single capitalized letter
}
