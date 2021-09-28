package handler

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	user_model "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/middleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserHandler struct {
	storage storage_user.UserUseCase
}

func NewLoginHandler(storageUser *storage_user.UserUseCase) *UserHandler {
	return &UserHandler{
		storage: *storageUser,
	}
}

func (handler *UserHandler) Login(ctx echo.Context) error {
	var newUserInput user_model.UserInput
	if err := ctx.Bind(&newUserInput); err != nil {
		newLoginError := user_model.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}
	fmt.Println(newUserInput)
	//if err := ctx.Validate(&newUserInput); err != nil {
	//	newLoginError := user_model.NewError(21, "validation error")
	//	return ctx.JSON(http.StatusBadRequest, newLoginError)
	//}

	id, err := handler.storage.IsUserExists(newUserInput)

	if err != nil {
		newLoginError := user_model.NewError(id, err.Error())
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	if id == errors.NO_USER_ERROR {
		newLoginError := user_model.NewError(errors.NO_USER_ERROR, errors.NO_USER_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newLoginError)
	}

	if id == errors.WRONG_PASSWORD_ERROR {
		newLoginError := user_model.NewError(errors.WRONG_PASSWORD_ERROR, errors.WRONG_PASSWORD_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newLoginError)
	}

	claimsString, err := middleware.GetToken(id, newUserInput.Name)
	if err != nil {
		newLoginError := user_model.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	handler.setCookieValue(ctx, claimsString)
	return ctx.JSON(http.StatusOK, newUserInput)
}

func (handler *UserHandler) SignUp(ctx echo.Context) error {
	var newUser user_model.User
	if err := ctx.Bind(&newUser); err != nil {
		newSignupError := user_model.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if err := ctx.Validate(&newUser); err != nil {
		newSignupError := user_model.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	newId, err := handler.storage.AddUser(newUser)
	if err != nil {
		newSignupError := user_model.NewError(newId, err.Error())
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if newId == errors.USER_EXISTS_ERROR {
		newSignupError := user_model.NewError(errors.USER_EXISTS_ERROR, errors.USER_EXISTS_DESCR)
		return ctx.JSON(http.StatusUnauthorized, newSignupError)
	}

	claimsString, err := middleware.GetToken(newId, newUser.Name)
	if err != nil {
		newSignupError := user_model.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	handler.setCookieValue(ctx, claimsString)
	return ctx.JSON(http.StatusOK, newUser)
}

func (handler *UserHandler) Logout(ctx echo.Context) error {
	cookie := &http.Cookie{
		Name:     "session_id",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
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
		Secure:   true,
	}

	ctx.SetCookie(cookie)
}
