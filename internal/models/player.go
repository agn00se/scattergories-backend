package models

// Using a composite key (UserID and GameID) here
// gorm.Model is omitted because it includes an ID field as the primary key, which we do not need.
type Player struct {
	UserID uint `gorm:"primaryKey;not null" json:"user_id"`
	GameID uint `gorm:"primaryKey;not null" json:"game_id"`
	Score  int  `gorm:"default:0" json:"score"`

	// User   User     `gorm:"foreignKey:UserID" json:"-"`
	// Game   Game     `gorm:"foreignKey:GameID" json:"-"`
}
