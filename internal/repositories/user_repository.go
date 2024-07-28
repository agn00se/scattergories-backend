package repositories

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUsersByGameRoomID(roomID uuid.UUID) ([]*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUserByID(id uuid.UUID) *gorm.DB
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUsersByGameRoomID(roomID uuid.UUID) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Where("game_room_id = ?", roomID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) DeleteUserByID(id uuid.UUID) *gorm.DB {
	result := r.db.Unscoped().Delete(&domain.User{}, id)
	if result.Error != nil {
		return result
	}
	if result.RowsAffected == 0 {
		result.Error = common.ErrUserNotFound
	}
	return result
}
