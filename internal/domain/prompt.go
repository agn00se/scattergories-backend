package domain

import "gorm.io/gorm"

type Prompt struct {
	gorm.Model
	Text string `gorm:"not null;unique" json:"text"`
}
