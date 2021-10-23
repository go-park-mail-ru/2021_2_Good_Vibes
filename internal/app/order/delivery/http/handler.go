package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type OrderHandler struct {
	useCase order.UseCase
}

func NewOrderHandler(useCase order.UseCase) *OrderHandler {
	return &OrderHandler{
		useCase: useCase,
	}
}

const trace = "OrderHandler"

func (oh *OrderHandler) PutOrder(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " PutOrder")

	var newOrder models.Order
	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	idString := claims["id"].(string)
	userId, err := strconv.ParseInt(idString, 10, 64)

	newOrder.UserId = int(userId)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	if err := ctx.Bind(&newOrder); err != nil {
		logger.Error(err)
		newOrderError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	if err := ctx.Validate(&newOrder); err != nil {
		logger.Error(err, newOrder)
		newOrderError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newOrderError)
	}

	orderId, err := oh.useCase.PutOrder(newOrder)
	if err != nil {
		logger.Error(err, newOrder)
		newOrderError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newOrderError)
	}

	newOrder.OrderId = orderId

	logger.Trace(trace + " success PutOrder")
	return ctx.JSON(http.StatusOK, newOrder)
}
