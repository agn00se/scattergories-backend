package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"

	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

func GetAllUsers() ([]*domain.User, error) {
	return repositories.GetAllUsers()
}

func GetUserByID(id uint) (*domain.User, error) {
	return repositories.GetUserByID(id)
}

func CreateGuestUser() (*domain.User, error) {
	guestName := utils.GenerateGuestName()

	user := &domain.User{
		Type: domain.UserTypeGuest,
		Name: guestName,
	}
	return createUser(user)
}

func CreateRegisteredUser(name string, email string, password string) (*domain.User, error) {
	hash, salt, err := generateHash(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Type:         domain.UserTypeRegistered,
		Name:         name,
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

func getUserByEmail(email string) (*domain.User, error) {
	return repositories.GetUserByEmail(email)
}

func getUsersByGameRoomID(roomID uint) ([]*domain.User, error) {
	return repositories.GetUsersByGameRoomID(roomID)
}

func createUser(user *domain.User) (*domain.User, error) {
	if err := repositories.CreateUser(user); err != nil {
		// Return error if the email is already used
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode {
			return nil, common.ErrEmailAlreadyUsed
		}
		return nil, err
	}
	return user, nil
}

func updateUser(user *domain.User) error {
	return repositories.UpdateUser(user)
}
