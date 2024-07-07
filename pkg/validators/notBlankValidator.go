package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom validation function to check for non-blank strings
func NotBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}
