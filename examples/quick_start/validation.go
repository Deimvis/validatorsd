package hi

import (
    "github.com/go-playground/validator/v10"

    "github.com/Deimvis/validatorsd"
)

var val = validator.New(validator.WithRequiredStructEnabled())

func Validate(obj interface{}) error {
    return validatorsd.Validate(val, obj)
}
