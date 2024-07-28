package domain

import (
	"time"

	"github.com/google/uuid"
)

type GameStatus string

const (
	GameStatusOngoing   GameStatus = "ongoing"
	GameStatusVoting    GameStatus = "voting"
	GameStatusCompleted GameStatus = "completed"
)

type Game struct {
	BaseModel
	GameRoomID  uuid.UUID    `gorm:"type:uuid;not null;index" json:"room_id"`
	Status      GameStatus   `gorm:"type:varchar(20);not null" json:"status"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	Players     []Player     `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;" json:"-"`
	GamePrompts []GamePrompt `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE;" json:"-"`
}
