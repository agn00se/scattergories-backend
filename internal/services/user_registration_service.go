package services

import (
	"scattergories-backend/internal/domain"
)

type UserRegistrationService interface {
	CreateRegisteredUser(name string, email string, password string) (*domain.User, error)
}

type UserRegistrationServiceImpl struct {
	userService UserService
	authService AuthService
}

func NewUserRegistrationService(userService UserService, authService AuthService) UserRegistrationService {
	return &UserRegistrationServiceImpl{userService: userService, authService: authService}
}

func (s *UserRegistrationServiceImpl) CreateRegisteredUser(name string, email string, password string) (*domain.User, error) {
	hash, salt, err := s.authService.GenerateHash(password)
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
	return s.userService.CreateUser(user)
}
