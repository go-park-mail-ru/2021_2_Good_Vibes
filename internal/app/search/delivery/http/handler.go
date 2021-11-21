package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SearchHandler struct {
	useCase        search.UseCase
}

func NewSearchHandler(useCase search.UseCase) *SearchHandler {
	return &SearchHandler{
		useCase:        useCase,
	}
}

const trace = "SearchHandler"

func (sh *SearchHandler) GetSuggests(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".GetSuggest")

	searchString := ctx.QueryParam("str")
	if searchString == "" {
		logger.Error("bad query param for GetSuggests")
		newError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	suggests, err := sh.useCase.GetSuggests(searchString)
	if err != nil {
		logger.Error(err)
		newError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newError)
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

	suggests, err := sh.useCase.GetSearchResults(searchArray)
	if err != nil {
		logger.Error(err)
		newError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	logger.Trace(trace + " success GetSearchResults")
	return ctx.JSON(http.StatusOK, suggests)
}