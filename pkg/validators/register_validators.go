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
		v.RegisterValidation("password_policy", PasswordValidator)
	}
}
