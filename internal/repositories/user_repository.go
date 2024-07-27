package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

func GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func GetUsersByGameRoomID(roomID uint) ([]*domain.User, error) {
	var users []*domain.User
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *domain.User) error {
	return config.DB.Create(user).Error
}

func UpdateUser(user *domain.User) error {
	return config.DB.Save(user).Error
}

func DeleteUserByID(id uint) *gorm.DB {
	result := config.DB.Unscoped().Delete(&domain.User{}, id)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrUserNotFound
	}
	return result
}
