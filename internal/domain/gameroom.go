package domain

import "github.com/google/uuid"

// struct tags

// hostID uuid: This field stores the foreign key reference to the User model's ID.
// It is used to establish the relationship at the database level.

// Host User: This field provides a way to represent the related User object in your application.
// It allows you to load the full User object when needed, rather than just the foreign key.

type GameRoom struct {
	BaseModel
	RoomCode       string         `gorm:"not null;unique" json:"room_code"`
	IsPrivate      bool           `gorm:"default:false" json:"is_private"`
	Passcode       string         `json:"passcode,omitempty"`         // Omits empty passcode field in JSON response
	HostID         uuid.UUID      `gorm:"type:uuid" json:"host_id"`   // HostID is automatically recognized as the foreign key for the Host field because it follows the naming convention FieldNameID (where FieldName is Host).
	Host           User           `gorm:"foreignKey:HostID" json:"-"` // Tells the JSON marshaller to ignore when serializing the struct to JSON.
	Games          []Game         `gorm:"foreignKey:GameRoomID;constraint:OnDelete:CASCADE;" json:"-"`
	GameRoomConfig GameRoomConfig `gorm:"foreignKey:GameRoomID;constraint:OnDelete:CASCADE;" json:"-"`
}
