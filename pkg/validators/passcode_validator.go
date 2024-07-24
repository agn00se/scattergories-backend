package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func PasscodeRequiredIfPrivate(fl validator.FieldLevel) bool {
	isPrivate := fl.Parent().FieldByName("IsPrivate").Bool()
	passcode := fl.Field().String()

	if isPrivate && strings.TrimSpace(passcode) == "" {
		return false
	}
	return true
}
