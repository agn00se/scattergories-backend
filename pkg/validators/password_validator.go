package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasDigit   = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[@$!%*?&]`)
	validChars = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
)

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 8 &&
		validChars.MatchString(password) &&
		hasLower.MatchString(password) &&
		hasUpper.MatchString(password) &&
		hasDigit.MatchString(password) &&
		hasSpecial.MatchString(password)
}
