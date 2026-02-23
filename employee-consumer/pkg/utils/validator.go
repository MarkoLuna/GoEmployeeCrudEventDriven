package utils

import (
	"github.com/MarkoLuna/EmployeeConsumer/internal/constants"
	"gopkg.in/go-playground/validator.v9"
)

func CreateValidator() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("EmployeeStatusValid", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field().String()
		return fieldValue == string(constants.ACTIVE) || fieldValue == string(constants.INACTIVE)
	})

	return v
}
