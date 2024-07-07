package models

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	Answer       string `gorm:"not null;size:30" json:"answer"`
	IsValid      bool   `gorm:"not null;default:false" json:"is_valid"`
	PlayerID     uint   `gorm:"not null" json:"player_id"`
	GamePromptID uint   `gorm:"not null" json:"game_prompt_id"`
}
