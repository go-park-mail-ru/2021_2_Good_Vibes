package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order"
	sessionJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/parser"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/sanitizer"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

const NewOrder = "новый"

type OrderHandler struct {
	useCase        order.UseCase
	sessionManager sessionJwt.TokenManager
}

func NewOrderHandler(useCase order.UseCase, sessionManager sessionJwt.TokenManager) *OrderHandler {
	return &OrderHandler{
		useCase:        useCase,
		sessionManager: sessionManager,
	}
}

const trace = "OrderHandler"

func (oh *OrderHandler) PutOrder(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " PutOrder")

	var newOrder models.Order
	userId, err := oh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	newOrder.UserId = int(userId)

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

	newOrder = sanitizer.SanitizeData(&newOrder).(models.Order)

	newOrder.Status = NewOrder

	if ctx.Request().RequestURI == "/api/cart/confirm" {
		orderId, orderCost, err := oh.useCase.PutOrder(newOrder)
		if err != nil {
			logger.Error(err, newOrder)
			newOrderError := errors.NewError(errors.SERVER_ERROR, errors.SERVER_ERROR_DESCR)
			return ctx.JSON(http.StatusInternalServerError, newOrderError)
		}

		newOrder.OrderId = orderId
		newOrder.Cost = orderCost
	}

	if ctx.Request().RequestURI == "/api/cart/check" {
		getOrder, err := oh.useCase.GetOrderPriceWithPromo(newOrder)
		if err != nil {
			logger.Error(err)
			newOrderError := errors.NewError(errors.SERVER_ERROR, "Ошибка в подтверждении заказа")
			return ctx.JSON(http.StatusInternalServerError, newOrderError)
		}
		newOrder = *getOrder
	}

	logger.Trace(trace + " success PutOrder")
	return ctx.JSON(http.StatusOK, newOrder)
}

func (oh *OrderHandler) GetAllOrders(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetAllOrders")

	userId, err := oh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	orders, err := oh.useCase.GetAllOrders(int(userId))
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if orders == nil {
		orders = make([]models.Order, 0)
	}

	for index, _ := range orders {
		orders[index].Date = parser.ParseDateFromSql2(orders[index].Date)
	}

	for i, _ := range orders {
		for j, _ := range orders[i].Products {
			imageSlice := strings.Split(orders[i].Products[j].Image, ";")
			orders[i].Products[j].Image = imageSlice[0]
		}
	}

	return ctx.JSON(http.StatusOK, orders)
}
