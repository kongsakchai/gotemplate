package validator

import (
	"github.com/go-playground/validator/v10"
)

type reqValidator struct {
	validator *validator.Validate
}

func NewReqValidator() *reqValidator {
	validate := validator.New()
	return &reqValidator{
		validator: validate,
	}
}

func (v *reqValidator) Validate(obj any) error {
	if err := v.validator.Struct(obj); err != nil {
		return err
	}
	return nil
}
