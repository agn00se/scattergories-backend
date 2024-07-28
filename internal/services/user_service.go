package services

import (
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"
	"scattergories-backend/internal/repositories"
	"scattergories-backend/pkg/utils"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const uniqueViolationCode = "23505"

type UserService interface {
	GetAllUsers() ([]*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	CreateGuestUser() (*domain.User, error)
	DeleteUserByID(id uuid.UUID) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUsersByGameRoomID(roomID uuid.UUID) ([]*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepository: userRepository}
}

func (s *UserServiceImpl) GetAllUsers() ([]*domain.User, error) {
	return s.userRepository.GetAllUsers()
}

func (s *UserServiceImpl) GetUserByID(id uuid.UUID) (*domain.User, error) {
	return s.userRepository.GetUserByID(id)
}

func (s *UserServiceImpl) CreateGuestUser() (*domain.User, error) {
	guestName := utils.GenerateGuestName()

	user := &domain.User{
		Type: domain.UserTypeGuest,
		Name: guestName,
	}
	return s.CreateUser(user)
}

func (s *UserServiceImpl) DeleteUserByID(id uuid.UUID) error {
	result := s.userRepository.DeleteUserByID(id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserServiceImpl) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepository.GetUserByEmail(email)
}

func (s *UserServiceImpl) GetUsersByGameRoomID(roomID uuid.UUID) ([]*domain.User, error) {
	return s.userRepository.GetUsersByGameRoomID(roomID)
}

func (s *UserServiceImpl) CreateUser(user *domain.User) (*domain.User, error) {
	if err := s.userRepository.CreateUser(user); err != nil {
		// Return error if the email is already used
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationCode {
			return nil, common.ErrEmailAlreadyUsed
		}
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) UpdateUser(user *domain.User) error {
	return s.userRepository.UpdateUser(user)
}
