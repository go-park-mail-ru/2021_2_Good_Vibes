package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type BasketHandler struct {
	useCase basket.UseCase
}

func NewBasketHandler(useCase basket.UseCase) *BasketHandler {
	return &BasketHandler{
		useCase: useCase,
	}
}

func (bh *BasketHandler) PutInBasket(ctx echo.Context) error {
	var newProduct models.BasketProduct
	if err := ctx.Bind(&newProduct); err != nil {
		newOrderError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	if err := ctx.Validate(&newProduct); err != nil {
		newOrderError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	err := bh.useCase.PutInBasket(newProduct)
	if err != nil {
		newOrderError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	return ctx.JSON(http.StatusOK, newProduct)
}

func (bh *BasketHandler) DropBasket(ctx echo.Context) error {
	type User struct {
		UserId int `json:"user_id"`
	}
	var user User
	if err := ctx.Bind(&user); err != nil {
		newOrderError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	if err := ctx.Validate(&user); err != nil {
		newOrderError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	err := bh.useCase.DropBasket(user.UserId)
	if err != nil {
		newOrderError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	return ctx.JSON(http.StatusOK, user)
}

func (bh *BasketHandler) DeleteProduct(ctx echo.Context) error {
	var product models.BasketProduct
	if err := ctx.Bind(&product); err != nil {
		newOrderError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	if err := ctx.Validate(&product); err != nil {
		newOrderError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	err := bh.useCase.DeleteProduct(product)
	if err != nil {
		newOrderError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	return ctx.JSON(http.StatusOK, product)
}

