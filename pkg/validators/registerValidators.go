package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("not_blank", NotBlank)
		v.RegisterValidation("passcode_required_if_private", PasscodeRequiredIfPrivate)
	}
}
