package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	middlewareLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	usecase product.Usecase
}

func NewProductHandler(usecase product.Usecase) *ProductHandler {
	return &ProductHandler{
		usecase: usecase,
	}
}

const trace = "ProductHandler"
//это должно быть, пока не нужно
/*
func addProduct(ctx echo.Context) error {
	return nil
}*/

//пока решаем вопросы с пагинацией - так
func (ph *ProductHandler) GetAllProducts(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + "GetAllProducts")

	answer, err := ph.usecase.GetAllProducts()
	if err != nil {
		logger.Error(err)
		newProductError := errors.NewError(errors.DB_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}
	return ctx.JSON(http.StatusOK, answer)
}

func (ph *ProductHandler) GetProductById(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + "GetProductById")

	val := ctx.QueryParams()

	idString := val.Get("id")
	if idString == "" {
		logger.Error("bad query param for GetProductById")
		newProductError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Error(err)
		newProductError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}

	answer, err := ph.usecase.GetProductById(id)
	if err != nil {
		logger.Error(err)
		newProductError := errors.NewError(errors.DB_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}

	return ctx.JSON(http.StatusOK, answer)
}
