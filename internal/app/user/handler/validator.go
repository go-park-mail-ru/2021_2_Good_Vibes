package handler

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
	"unicode"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func Password(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	number, upper, special, letter := false, false, false, false
	for _, c := range val {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letter = true
		default:
			return false
		}
	}
	return number && upper && special && letter && len(val) >= 7 && len(val) <= 20
}
