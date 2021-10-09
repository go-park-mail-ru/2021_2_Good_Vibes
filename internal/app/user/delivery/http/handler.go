package http

import (
	"github.com/dgrijalva/jwt-go"
	errors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	session "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	Usecase user.Usecase
}

func NewLoginHandler(storageUser user.Usecase) *UserHandler {
	return &UserHandler{
		Usecase: storageUser,
	}
}

func (handler *UserHandler) Login(ctx echo.Context) error {
	var newUserDataForInput models.UserDataForInput
	if err := ctx.Bind(&newUserDataForInput); err != nil {
		newLoginError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	if err := ctx.Validate(&newUserDataForInput); err != nil {
		newLoginError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	id, err := handler.Usecase.CheckPassword(newUserDataForInput)

	if err != nil {
		newLoginError := errors.NewError(id, err.Error())
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	if id == errors.NO_USER_ERROR {
		newLoginError := errors.NewError(errors.NO_USER_ERROR, errors.NO_USER_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newLoginError)
	}

	if id == errors.WRONG_PASSWORD_ERROR {
		newLoginError := errors.NewError(errors.WRONG_PASSWORD_ERROR, errors.WRONG_PASSWORD_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newLoginError)
	}

	claimsString, err := session.GetToken(id, newUserDataForInput.Name)
	if err != nil {
		newLoginError := errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	handler.setCookieValue(ctx, claimsString)
	newUserDataForInput.Password = ""
	return ctx.JSON(http.StatusOK, newUserDataForInput)
}

func (handler *UserHandler) SignUp(ctx echo.Context) error {
	var newUser models.UserDataForReg
	if err := ctx.Bind(&newUser); err != nil {
		newSignupError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if err := ctx.Validate(&newUser); err != nil {
		newSignupError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	newId, err := handler.Usecase.AddUser(newUser)
	if err != nil {
		newSignupError := errors.NewError(newId, err.Error())
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if newId == errors.USER_EXISTS_ERROR {
		newSignupError := errors.NewError(errors.USER_EXISTS_ERROR, errors.USER_EXISTS_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newSignupError)
	}

	claimsString, err := session.GetToken(newId, newUser.Name)
	if err != nil {
		newSignupError := errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	handler.setCookieValue(ctx, claimsString)
	newUser.Password = ""
	return ctx.JSON(http.StatusOK, newUser)
}

func (handler *UserHandler) Profile(ctx echo.Context) error {
	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	idString := claims["id"].(string)
	idNum, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	userData, err := handler.Usecase.GetUserDataByID(idNum)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, err.Error()))
	}

	return ctx.JSON(http.StatusOK, userData)
}

func (handler *UserHandler) Logout(ctx echo.Context) error {
	cookie := &http.Cookie{
		Name:     "session_id",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		//Secure:   true,
	}
	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (handler *UserHandler) setCookieValue(ctx echo.Context, value string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    value,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 72),
		SameSite: http.SameSiteNoneMode,
		//Secure:   true,
	}

	ctx.SetCookie(cookie)
}
