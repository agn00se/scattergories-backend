package models

import "gorm.io/gorm"

// json tags here are not strictly necessary as we use toUserResponse to convert into response struct.
// But keeping for various benefits:
// Direct Serialization: If you ever need to serialize the model directly to JSON (e.g., for debugging or logging purposes).
// Consistency: Maintaining consistent use of json tags across your codebase can help avoid errors and make the code more understandable.
type User struct {
	gorm.Model
	Name       string `gorm:"not null" json:"name"`
	GameRoomID *uint  `gorm:"index" json:"room_id,omitempty"` // Nullable and indexed for efficient querying
}
