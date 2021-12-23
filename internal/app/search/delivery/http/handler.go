package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SearchHandler struct {
	useCase search.UseCase
}

func NewSearchHandler(useCase search.UseCase) *SearchHandler {
	return &SearchHandler{
		useCase: useCase,
	}
}

const trace = "SearchHandler"

func (sh *SearchHandler) GetSuggests(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".GetSuggest")

	searchString := ctx.QueryParam("str")

	suggests, err := sh.useCase.GetSuggests(searchString)
	if err != nil {
		logger.Error(err)
		newError := errors.NewError(errors.SERVER_ERROR, errors.SERVER_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if suggests.Products == nil {
		suggests.Products = make([]models.ProductForSuggest, 0)
	}

	if suggests.Categories == nil {
		suggests.Categories = make([]models.CategoryForSuggest, 0)
	}

	logger.Trace(trace + " success GetSuggest")
	return ctx.JSON(http.StatusOK, suggests)
}

func (sh *SearchHandler) GetSearchResults(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".GetSearchResults")

	searchString := ctx.QueryParam("str")
	if searchString == "" {
		logger.Error("bad query param for GetSearchResults")
		newError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	searchArray := strings.Split(searchString, " ")

	filter, err := postgre.ParseQueryFilter(ctx)
	if err != nil {
		logger.Trace(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	suggests, err := sh.useCase.GetSearchResults(searchArray, *filter)
	if err != nil {
		logger.Error(err)
		newError := errors.NewError(errors.SERVER_ERROR, errors.SERVER_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if suggests.Products == nil {
		suggests.Products = make([]models.Product, 0)
	}

	for i, _ := range suggests.Products {
		imageSlice := strings.Split(suggests.Products[i].Image, ";")
		suggests.Products[i].Image = imageSlice[0]
	}

	logger.Trace(trace + " success GetSearchResults")
	return ctx.JSON(http.StatusOK, suggests)
}
