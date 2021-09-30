package configValidator

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/handler"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ConfigValidator(router *echo.Echo) {
	val := validator.New()
	val.RegisterValidation("customPassword", handler.Password)
	router.Validator = &handler.CustomValidator{Validator: val}
}
