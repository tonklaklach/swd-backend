package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/tonklaklach/swd-backend/4_rest_api_go/domain"
)

var requestValidator = newRequestValidator()

func newRequestValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		return strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	})
	if err := v.RegisterValidation("notblank", func(field validator.FieldLevel) bool {
		return strings.TrimSpace(field.Field().String()) != ""
	}); err != nil {
		panic(err)
	}
	return v
}

func Validate(req any) error {
	err := requestValidator.Struct(req)
	if err == nil {
		return nil
	}
	var fields validator.ValidationErrors
	if !errors.As(err, &fields) {
		return err
	}
	field := fields[0]
	var message string
	switch field.Tag() {
	case "required", "notblank":
		message = fmt.Sprintf("%s is required", field.Field())
	case "gt":
		message = fmt.Sprintf("%s must be positive", field.Field())
	case "oneof":
		message = fmt.Sprintf("%s must be one of %s", field.Field(), field.Param())
	case "unique":
		message = fmt.Sprintf("%s must contain unique values", field.Field())
	default:
		message = fmt.Sprintf("%s is invalid", field.Field())
	}
	return &domain.ValidationError{Message: message}
}
