package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

// newValidator membuat validator yang melaporkan error memakai nama field JSON
// (mis. "kode_barang"), bukan nama field struct Go (mis. "KodeBarang").
func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return strings.ToLower(f.Name)
		}
		return name
	})
	return v
}

func ValidateStruct(payload any) map[string]string {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrs {
			field := fe.Field()
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
	case "gte":
		return fmt.Sprintf("Field %s minimal %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("Field %s maksimal %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("Field %s harus lebih besar dari %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("Field %s harus lebih kecil dari %s", field, fe.Param())
	case "datetime":
		return fmt.Sprintf("Field %s harus berformat tanggal YYYY-MM-DD", field)
	case "oneof":
		return fmt.Sprintf("Field %s harus salah satu dari: %s", field, fe.Param())
	default:
		return fmt.Sprintf("Field %s tidak valid (%s)", field, fe.Tag())
	}
}
