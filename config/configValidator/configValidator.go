package configValidator

import (
	validator2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/validator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ConfigValidator(router *echo.Echo) {
	val := validator.New()
	val.RegisterValidation("customPassword", validator2.Password)
	router.Validator = &validator2.CustomValidator{Validator: val}
}
