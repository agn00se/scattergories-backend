package models

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	GameRoomID uint       `gorm:"not null" json:"room_id"`
	Status     GameStatus `gorm:"type:varchar(20);not null" json:"status"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    time.Time  `json:"end_time"`
	Letter     string     `gorm:"type:char(1);not null" json:"letter"` // Store as a single capitalized letter

	// Room      GameRoom   `gorm:"foreignKey:RoomID" json:"-"`
}
