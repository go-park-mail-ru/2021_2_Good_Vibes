package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	middlewareLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

var AllCategoriesJson models.CategoryNode

type CategoryHandler struct {
	useCase category.UseCase
}

func NewCategoryHandler(useCase category.UseCase) *CategoryHandler {
	return &CategoryHandler{
		useCase: useCase,
	}
}

const trace = "CategoryHandler"

func (ch *CategoryHandler) GetCategories(ctx echo.Context) error {
	logger := ctx.Get(middlewareLogger.LoggerFieldName).(*logrus.Entry)
	logger.Trace(trace + " GetCategories")

	val := ctx.QueryParams()

	nameString := val.Get("name")
	if nameString == "" {
		if AllCategoriesJson.Name != "" {
			logger.Debug(AllCategoriesJson)
			return ctx.JSON(http.StatusOK, AllCategoriesJson)
		}

		categories, err := ch.useCase.GetAllCategories()
		if err != nil {
			logger.Error(err)
			newCategoryError := errors.NewError(errors.DB_ERROR, err.Error())
			return ctx.JSON(http.StatusBadRequest, newCategoryError)
		}
		AllCategoriesJson = categories

		logger.Debug(AllCategoriesJson)
		return ctx.JSON(http.StatusOK, categories)
	}

	products, err := ch.useCase.GetProductsByCategory(nameString)
	if err != nil {
		logger.Error(err)
		newCategoryError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	logger.Debug(products)
	return ctx.JSON(http.StatusOK, products)
}
