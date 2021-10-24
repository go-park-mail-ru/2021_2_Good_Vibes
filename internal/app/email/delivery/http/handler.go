package http

import (
	emailPac "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/email"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EmailHandler struct {
	useCase emailPac.UseCase
}

func NewEmailHandler(useCase emailPac.UseCase) *EmailHandler {
	return &EmailHandler{
		useCase: useCase,
	}
}

func (ch *EmailHandler) ConfirmEmail(ctx echo.Context) error {
	val := ctx.QueryParams()
	email := val.Get("email")
	token := val.Get("token")

	err := ch.useCase.ConfirmEmail(email, token)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	message := "Registration confirmed!"

	return ctx.JSON(http.StatusOK, message)
}
