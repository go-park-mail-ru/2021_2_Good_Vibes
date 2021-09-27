package handler

import (
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
	newUserInput := new(user_model.UserInput)
	if err := ctx.Bind(newUserInput); err != nil {
		newLoginError := user_model.NewError(20, "cannot bind data")
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}
	if err := ctx.Validate(newUserInput); err != nil {
		newLoginError := user_model.NewError(21, "validation error")
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	id, err := handler.storage.IsUserExists(*newUserInput)
	if id == -1 {
		if err != nil {
			newLoginError := user_model.NewError(31, err.Error())
			return ctx.JSON(http.StatusUnauthorized, newLoginError)
		}
		newLoginError := user_model.NewError(30, "user does not exist")
		return ctx.JSON(http.StatusUnauthorized, newLoginError)
	}

	claimsString, err := middleware.GetToken(id, newUserInput.Name)
	if err != nil {
		newLoginError := user_model.NewError(22, "cannot get token")
		return ctx.JSON(http.StatusBadRequest, newLoginError)
	}

	handler.setCookieValue(ctx, claimsString)
	return ctx.JSON(http.StatusOK, newUserInput)
}

func (handler *UserHandler) SignUp(ctx echo.Context) error {
	newUser := new(user_model.User)
	if err := ctx.Bind(newUser); err != nil {
		newSignupError := user_model.NewError(20, "cannot bind data")
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}
	if err := ctx.Validate(newUser); err != nil {
		newSignupError := user_model.NewError(21, "validation error")
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	newId, err := handler.storage.AddUser(*newUser)
	if err != nil {
		newSignupError := user_model.NewError(40, err.Error())
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if newId == -1 {
		newSignupError := user_model.NewError(32, "user exists")
		return ctx.JSON(http.StatusUnauthorized, newSignupError)
	}

	claimsString, err := middleware.GetToken(newId, newUser.Name)
	if err != nil {
		newSignupError := user_model.NewError(20, "cannot get token")
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
