package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func GetUsersByGameRoomID(roomID uint) ([]*models.User, error) {
	var users []*models.User
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func DeleteUserByID(id uint) *gorm.DB {
	result := config.DB.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrUserNotFound
	}
	return result
}
