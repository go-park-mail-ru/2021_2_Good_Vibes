package http

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	middlewareLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type BasketHandler struct {
	useCase basket.UseCase
}

func NewBasketHandler(useCase basket.UseCase) *BasketHandler {
	return &BasketHandler{
		useCase: useCase,
	}
}

const trace = "BasketHandler"

func (bh *BasketHandler) PutInBasket(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + ".PutInBasket")

	var newProduct models.BasketProduct
	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	idString := claims["id"].(string)
	idNum, err := strconv.ParseInt(idString, 10, 64)

	newProduct.UserId = int(idNum)
	if err != nil {
		logger.Error(err, newProduct.UserId)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	if err := ctx.Bind(&newProduct); err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	if err := ctx.Validate(&newProduct); err != nil {
		logger.Error(err, newProduct)
		newBasketError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	err = bh.useCase.PutInBasket(newProduct)
	if err != nil {
		logger.Error(err, newProduct)
		newBasketError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newBasketError)
	}

	logger.Trace(trace + " success PutInBasket")
	return ctx.JSON(http.StatusOK, newProduct)
}

func (bh *BasketHandler) GetBasket(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + ".GetBasket")

	var user models.UserID

	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	stringId := claims["id"].(string)

	userId, err := strconv.ParseInt(stringId, 10, 64)

	user.UserId = int(userId)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	var basketProducts []models.BasketProduct

	basketProducts, err = bh.useCase.GetBasket(user.UserId)
	if err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	logger.Trace(trace + " success GetBasket")
	return ctx.JSON(http.StatusOK, basketProducts)
}

func (bh *BasketHandler) DropBasket(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + ".DropBasket")

	var user models.UserID

	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	stringId := claims["id"].(string)

	userId, err := strconv.ParseInt(stringId, 10, 64)
	user.UserId = int(userId)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	if err := ctx.Bind(&user); err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	if err := ctx.Validate(&user); err != nil {
		logger.Error(err, user)
		newBasketError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	err = bh.useCase.DropBasket(user.UserId)
	if err != nil {
		logger.Error(err, user)
		newBasketError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	logger.Trace(trace + " success DropBasket")
	return ctx.JSON(http.StatusOK, user)
}

func (bh *BasketHandler) DeleteProduct(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + ".DeleteProduct")

	var product models.BasketProduct

	token := ctx.Get("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	stringId := claims["id"].(string)

	userId, err := strconv.ParseInt(stringId, 10, 64)

	product.UserId = int(userId)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	if err := ctx.Bind(&product); err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	if err := ctx.Validate(&product); err != nil {
		logger.Error(err, product)
		newBasketError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	err = bh.useCase.DeleteProduct(product)
	if err != nil {
		logger.Error(err, product)
		newBasketError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	logger.Trace(trace + " success DeleteProduct")
	return ctx.JSON(http.StatusOK, product)
}
