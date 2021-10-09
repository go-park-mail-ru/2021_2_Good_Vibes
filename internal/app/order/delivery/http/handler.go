package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order"
	"github.com/labstack/echo/v4"
	"net/http"
)

type OrderHandler struct {
	useCase order.UseCase
}

func NewOrderHandler(useCase order.UseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

func (oh *OrderHandler) PutOrder(ctx echo.Context) error {
	var newOrder models.Order
	if err := ctx.Bind(&newOrder); err != nil {
		newOrderError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	if err := ctx.Validate(&newOrder); err != nil {
		newOrderError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	orderId, err := oh.useCase.PutOrder(newOrder)
	if err != nil {
		newOrderError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	newOrder.OrderId = orderId
	return ctx.JSON(http.StatusOK, newOrder)
}
