package http

import (
	"crypto/tls"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	emailPac "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/email"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/labstack/echo/v4"
	gomail "gopkg.in/mail.v2"
	"net/http"
	"strconv"
)

type EmailHandler struct {
	useCase emailPac.UseCase
}

func NewEmailHandler(useCase emailPac.UseCase) *EmailHandler {
	return &EmailHandler{
		useCase: useCase,
	}
}

func (eh *EmailHandler) ConfirmEmail(ctx echo.Context) error {
	val := ctx.QueryParams()
	email := val.Get("email")
	token := val.Get("token")

	err := eh.useCase.ConfirmEmail(email, token)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	message := "Registration confirmed!"

	return ctx.JSON(http.StatusOK, message)
}

func (eh *EmailHandler) SendConfirmationEmail(ctx echo.Context) error {
	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	idString := claims["id"].(string)
	idNum, err := strconv.ParseInt(idString, 10, 64)

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	email, err := eh.useCase.GetUserEmailById(uint64(idNum))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	emailToken := eh.useCase.GenerateToken()

	err = eh.useCase.InsertUserToken(email, emailToken)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", config.ConfigApp.Email.Address)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Подтверждение регистрации")
	body := fmt.Sprintf(`Здравствуйте! Чтобы подтвердить регистрацию, кликните по ссылке 
                                <a href = "http://127.0.0.1:8080/email/confirm?email=%s&token=%s">AZOT</a>.`, email, emailToken)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer(config.ConfigApp.Email.Server,
		config.ConfigApp.Email.ServerPort,
		config.ConfigApp.Email.Address,
		config.ConfigApp.Email.Password)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(message); err != nil {
		return ctx.JSON(http.StatusBadRequest, errors.NewError(errors.SERVER_ERROR, "Email send error"))
	}

	return ctx.JSON(http.StatusOK, "Email sent")
}
