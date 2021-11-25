package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/sanitizer"
	"github.com/labstack/echo/v4"
	"net/http"
)

// var AllCategoriesJson models.CategoryNode

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
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetCategories")

	categories, err := ch.useCase.GetAllCategories()
	if err != nil {
		logger.Error(err)
		newCategoryError := errors.NewError(errors.DB_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newCategoryError)
	}

	categories = sanitizer.SanitizeData(&categories).(models.CategoryNode)

	return ctx.JSON(http.StatusOK, categories)
}

func (ch *CategoryHandler) GetCategoryProducts(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetCategoryProducts")

	filter, err := postgre.ParseQueryFilter(ctx)
	if err != nil {
		logger.Trace(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	nameString := ctx.Param("name")
	filter.NameCategory = nameString

	products, err := ch.useCase.GetProductsByCategory(*filter)
	if err != nil {
		logger.Error(err)
		newCategoryError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	for i, _ := range products {
		products[i] = sanitizer.SanitizeData(&products[i]).(models.Product)
	}

	if products == nil {
		products = make([]models.Product, 0)
	}

	logger.Debug(products)
	return ctx.JSON(http.StatusOK, products)
}

func (ch *CategoryHandler) CreateCategory(ctx echo.Context) error {
	var newCategory models.CreateCategory

	if err := ctx.Bind(&newCategory); err != nil {
		newSignupError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if err := ctx.Validate(&newCategory); err != nil {
		newCategoryError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	newCategory = sanitizer.SanitizeData(&newCategory).(models.CreateCategory)

	err := ch.useCase.CreateCategory(newCategory)
	if err != nil {
		newCategoryError := errors.NewError(errors.SERVER_ERROR, errors.BD_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	categories, err := ch.useCase.GetAllCategories()
	if err != nil {
		newCategoryError := errors.NewError(errors.SERVER_ERROR, errors.BD_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	categories = sanitizer.SanitizeData(&categories).(models.CategoryNode)

	// AllCategoriesJson = categories

	return ctx.JSON(http.StatusOK, categories)
}
