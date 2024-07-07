package models

import "gorm.io/gorm"

type GamePrompt struct {
	gorm.Model
	GameID uint   `gorm:"not null" json:"game_id"`
	Text   string `gorm:"not null" json:"text"`
}
