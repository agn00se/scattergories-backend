package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validate = v
		v.RegisterValidation("not_blank", NotBlank)
		v.RegisterValidation("passcode_required_if_private", PasscodeRequiredIfPrivate)
		v.RegisterValidation("is_valid_user_type", IsValidUserType)
		v.RegisterValidation("name_required_if_registered", NameRequiredIfRegistered)
		v.RegisterValidation("email_required_if_registered", EmailRequiredIfRegistered)
		v.RegisterValidation("password_required_if_registered", PasswordRequiredIfRegistered)
	}
}
