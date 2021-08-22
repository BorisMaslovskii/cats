// Package validator validates incoming fields
package validator

import (
	validator "github.com/go-playground/validator/v10"
)

// CustomValidator struct
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates new validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate function
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
