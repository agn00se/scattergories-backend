package validators

import (
	"scattergories-backend/internal/models"

	"github.com/go-playground/validator/v10"
)

func IsValidUserType(fl validator.FieldLevel) bool {
	userType := models.UserType(fl.Field().String())
	switch userType {
	case models.UserTypeRegistered, models.UserTypeGuest:
		return true
	}
	return false
}
