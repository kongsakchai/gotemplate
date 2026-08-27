package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: validate,
	}
}

func (v *Validator) Validate(obj any) error {
	if err := v.validator.Struct(obj); err != nil {
		errMap := make(errorMap)
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				fieldName := e.Field()
				tag := e.Tag()
				errMap[fieldName] = tag
			}
		}
		return errMap
	}
	return nil
}
