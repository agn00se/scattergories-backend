package services

import (
	"fmt"
	"math/rand"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
	"time"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

func GetAllUsers() ([]*models.User, error) {
	return repositories.GetAllUsers()
}

func GetUserByID(id uint) (*models.User, error) {
	return repositories.GetUserByID(id)
}

func CreateGuestUser() (*models.User, error) {
	guestName := generateRandomGuestName()

	user := &models.User{
		Name: guestName,
		Type: models.UserTypeGuest,
	}
	return createUser(user)
}

func CreateRegisteredUser(name string, email string, password string) (*models.User, error) {
	hash, salt, err := generateHash(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Type:         models.UserTypeGuest,
		Email:        &email,
		PasswordHash: &hash,
		Salt:         &salt,
	}
	return createUser(user)
}

func DeleteUserByID(id uint) error {
	result := repositories.DeleteUserByID(id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func generateRandomGuestName() string {
	return fmt.Sprintf("Guest%d", rand.Intn(10000))
}

func getUserByEmail(email string) (*models.User, error) {
	return repositories.GetUserByEmail(email)
}

func getUsersByGameRoomID(roomID uint) ([]*models.User, error) {
	return repositories.GetUsersByGameRoomID(roomID)
}

func createUser(user *models.User) (*models.User, error) {
	if err := repositories.CreateUser(user); err != nil {
		// Return error if the email is already used
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode {
			return nil, common.ErrEmailAlreadyUsed
		}
		return nil, err
	}
	return user, nil
}

func updateUser(user *models.User) error {
	return repositories.UpdateUser(user)
}
