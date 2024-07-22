package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	UserID  uint     `gorm:"not null;uniqueIndex:idx_game_user" json:"user_id"`
	GameID  uint     `gorm:"not null;uniqueIndex:idx_game_user" json:"game_id"`
	Score   int      `gorm:"default:0" json:"score"`
	User    User     `gorm:"foreignKey:UserID" json:"user"`
	Answers []Answer `gorm:"foreignKey:PlayerID;constraint:OnDelete:CASCADE;" json:"-"`
}
