package services

import (
	"fmt"
	"math/rand"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/repositories"
	"time"
)

func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

func generateRandomGuestName() string {
	return fmt.Sprintf("Guest%d", rand.Intn(10000))
}

func CreateGuestUser() (*models.User, error) {
	guestName := generateRandomGuestName()
	return CreateUser(guestName)
}

func CreateUser(name string) (*models.User, error) {
	user := &models.User{
		Name: name,
	}
	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]*models.User, error) {
	return repositories.GetAllUsers()
}

func GetUserByID(id uint) (*models.User, error) {
	return repositories.GetUserByID(id)
}

func GetUsersByGameRoomID(roomID uint) ([]*models.User, error) {
	return repositories.GetUsersByGameRoomID(roomID)
}

func UpdateUser(user *models.User) error {
	return repositories.UpdateUser(user)
}

func DeleteUserByID(id uint) error {
	result := repositories.DeleteUserByID(id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
