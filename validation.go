package validatorsd

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	ValidateSelf() error
}

// IsValid is a shortcut for Validate
// that returns whether err == nil.
func IsValid(v *validator.Validate, obj interface{}) bool {
	err := Validate(v, obj)
	return err == nil
}

// Validate validates given struct using both
// validation tags and ValidateSelf methods.
// It panics when given struct is nil.
// It goes recursively through substructs, including embedded structs.
// Doesn't matter whether ValidateSelf has a value or a pointer receiver -
// Validate function will find it either way.
func Validate(v *validator.Validate, obj interface{}) error {
	err := ValidateTags(v, obj)
	if err != nil {
		return err
	}
	err = ValidateSelfRecursively(obj)
	return err
}

func ValidateTags(v *validator.Validate, obj interface{}) error {
	return v.Struct(obj)
}

func ValidateSelfRecursively(obj interface{}) error {
	return validateSelfRecursively(reflect.ValueOf(obj), false)
}

func validateSelfRecursively(v reflect.Value, ignoreNil bool) error {
	return validate(extractInternalValue(v), ignoreNil)
}

func validate(v reflect.Value, ignoreNil bool) error {
	var err error
	switch v.Kind() {
	case reflect.Struct:
		err = validateStruct(v)
	case reflect.Slice, reflect.Array:
		err = validateSlice(v)
	case reflect.Ptr: // nil pointer
		if !v.IsNil() {
			panic("bug: validate was called on external value (non-nil pointer)")
		}
		if ignoreNil {
			err = nil
		} else {
			err = validateNilPtr(v)
		}
	default:
		err = nil
	}
	return err
}

func validateStruct(v reflect.Value) error {
	err := validateSelf(v)
	if err != nil {
		return err
	}
	for i := 0; i < v.NumField(); i++ {
		err = validateSelfRecursively(v.Field(i), true /*ignore nil*/)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateSlice(v reflect.Value) error {
	for i := 0; i < v.Len(); i++ {
		err := validateSelfRecursively(v.Index(i), true /*ignore nil*/)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateNilPtr(v reflect.Value) error {
	switch v.Type().Elem().Kind() {
	case reflect.Slice, reflect.Array:
		return nil
	}
	return errors.New("object is invalid: nil pointer")
}

func validateSelf(v reflect.Value) error {
	if isValidatable(v) {
		return validateSelfDo(v)
	} else if ptr := getPtr(v); isValidatable(ptr) {
		return validateSelfDo(ptr)
	}
	return nil
}

func isValidatable(v reflect.Value) bool {
	validatable := reflect.TypeOf((*Validatable)(nil)).Elem()
	isNil := (v.Kind() == reflect.Ptr && v.IsNil())
	return !isNil && v.CanInterface() && v.Type().Implements(validatable)
}

func validateSelfDo(v reflect.Value) error {
	validatable, ok := v.Interface().(Validatable)
	if !ok {
		panic("value doesn't implement Validatable")
	}
	return validatable.ValidateSelf()
}

func getPtr(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v.Addr()
	}
	ptr := reflect.New(v.Type())
	ptr.Elem().Set(v)
	return ptr
}

// https://github.com/go-playground/validator/blob/a947377040f8ebaee09f20d09a745ec369396793/util.go#L15
func extractInternalValue(current reflect.Value) reflect.Value {

BEGIN:
	switch current.Kind() {
	case reflect.Ptr:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		if current.IsNil() {
			return current
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current

	default:
		return current
	}
}
