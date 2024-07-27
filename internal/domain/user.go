package domain

import "gorm.io/gorm"

type UserType string

const (
	UserTypeGuest      UserType = "guest"
	UserTypeRegistered UserType = "registered"
)

// json tags here are not strictly necessary as we use ToUserResponse to convert into response struct.
// But keeping for various benefits:
// Direct Serialization: If ever needed to serialize the model directly to JSON (e.g., for debugging or logging purposes).
// Consistency: Maintaining consistent use of json tags across your codebase can help avoid errors and make the code more understandable.
type User struct {
	gorm.Model
	Type         UserType `gorm:"type:varchar(20);not null" json:"type"`
	Name         string   `gorm:"not null" json:"name"`
	Email        *string  `gorm:"uniqueIndex:idx_email,where:email IS NOT NULL" json:"email,omitempty"`
	PasswordHash *string  `json:"-"`
	Salt         *string  `json:"-"`
	GameRoomID   *uint    `gorm:"index" json:"room_id,omitempty"` // Nullable and indexed for efficient querying
}
