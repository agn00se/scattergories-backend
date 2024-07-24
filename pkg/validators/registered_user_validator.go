package validators

import (
	"scattergories-backend/internal/models"
	"strings"

	"github.com/go-playground/validator/v10"
)

func NameRequiredIfRegistered(fl validator.FieldLevel) bool {
	userType := fl.Parent().FieldByName("Type").String()
	name := fl.Field().String()

	if userType == string(models.UserTypeRegistered) && strings.TrimSpace(name) == "" {
		return false
	}
	return true
}

func EmailRequiredIfRegistered(fl validator.FieldLevel) bool {
	userType := fl.Parent().FieldByName("Type").String()
	email := fl.Field().String()

	if userType == string(models.UserTypeRegistered) && strings.TrimSpace(email) == "" {
		return false
	}
	return true
}

func PasswordRequiredIfRegistered(fl validator.FieldLevel) bool {
	userType := fl.Parent().FieldByName("Type").String()
	password := fl.Field().String()

	if userType == string(models.UserTypeRegistered) && strings.TrimSpace(password) == "" {
		return false
	}
	return true
}
