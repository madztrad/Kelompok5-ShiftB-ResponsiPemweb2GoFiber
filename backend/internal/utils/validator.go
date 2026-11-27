package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(payload any) map[string]string {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := strings.ToLower(fe.Field())
			errs[field] = formatMessage(field, fe)
		}
	}
	return errs
}

func formatMessage(field string, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Field %s wajib diisi", field)
	case "email":
		return fmt.Sprintf("Field %s harus berupa email yang valid", field)
	case "min":
		return fmt.Sprintf("Field %s minimal %s karakter/nilai", field, fe.Param())
	case "max":
		return fmt.Sprintf("Field %s maksimal %s karakter/nilai", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("Field %s harus salah satu dari: %s", field, fe.Param())
	default:
		return fmt.Sprintf("Field %s tidak valid (%s)", field, fe.Tag())
	}
}
