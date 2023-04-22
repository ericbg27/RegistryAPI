package api

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	phoneRegexp    = `^\d{8}$`
	passwordRegexp = `[^_#%]*`
)

var isPhone validator.Func = func(fl validator.FieldLevel) bool {
	phone, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	matched, err := regexp.Match(phoneRegexp, []byte(phone))
	if err != nil {
		return false
	}

	return matched
}

var validPassword validator.Func = func(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	matched, err := regexp.Match(passwordRegexp, []byte(password))
	if err != nil {
		return false
	}

	return matched
}
