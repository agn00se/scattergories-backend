package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	GameRoomID  uint         `gorm:"not null" json:"room_id"`
	Status      GameStatus   `gorm:"type:varchar(20);not null" json:"status"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	Players     []Player     `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;" json:"-"`
	GamePrompts []GamePrompt `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;" json:"-"`
}
