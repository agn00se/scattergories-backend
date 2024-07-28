package domain

type Prompt struct {
	BaseModel
	Text string `gorm:"not null;unique" json:"text"`
}
