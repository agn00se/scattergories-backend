package services

import (
	"fmt"
	"math/rand"
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

func init() {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

func generateRandomGuestName() string {
	return fmt.Sprintf("Guest%d", rand.Intn(10000))
}

func CreateGuestUser() (models.User, error) {
	guestName := generateRandomGuestName()
	return CreateUser(guestName)
}

func CreateUser(name string) (models.User, error) {
	user := models.User{Name: name}
	if err := config.DB.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func GetUserByID(id uint) (models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

func UpdateUserByID(id uint, newName string, newGameRoomID *uint) (models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, ErrUserNotFound
		}
		return user, err
	}
	user.Name = newName
	user.GameRoomID = newGameRoomID
	if err := config.DB.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func DeleteUserByID(id uint) error {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		return err
	}

	if err := config.DB.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
