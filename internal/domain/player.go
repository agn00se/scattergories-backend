package domain

import "github.com/google/uuid"

type Player struct {
	BaseModel
	UserID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_game_user" json:"user_id"`
	GameID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_game_user" json:"game_id"`
	Score   int       `gorm:"default:0" json:"score"`
	User    User      `gorm:"foreignKey:UserID" json:"user"`
	Answers []Answer  `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE;" json:"-"`
}
