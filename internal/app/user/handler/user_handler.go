package handler

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"github.com/labstack/echo/v4"
	"net/http"
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
	newUserInput := new(storage_user.UserInput)
	if err := ctx.Bind(newUserInput); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	if err := ctx.Validate(newUserInput); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	fmt.Println(handler.storage)
	id, err := handler.storage.IsUserExists(*newUserInput)
	if id == -1 || err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	fmt.Println(id)
	return ctx.JSON(http.StatusOK, newUserInput)
}

func (handler *UserHandler) SignUp(ctx echo.Context) error {
	newUser := new(storage_user.User)
	if err := ctx.Bind(newUser); err != nil {
		fmt.Println(err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	if err := ctx.Validate(newUser); err != nil {
		fmt.Println(err)
		return ctx.NoContent(http.StatusBadRequest)
	}
	fmt.Println(handler.storage)
	newId, _ := handler.storage.AddUser(*newUser)
	fmt.Println(newId)
	if newId == -1 {
		return ctx.JSON(http.StatusBadRequest, newUser)
	}
	return ctx.JSON(http.StatusOK, newUser)
}
