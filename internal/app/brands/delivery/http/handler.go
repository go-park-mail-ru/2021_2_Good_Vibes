package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	sessionJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/sanitizer"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type BrandHandler struct {
	useCase        brands.UseCase
	sessionManager sessionJwt.TokenManager
}

func NewBrandHandler(useCase brands.UseCase, sessionManager sessionJwt.TokenManager) *BrandHandler {
	return &BrandHandler{
		useCase:        useCase,
		sessionManager: sessionManager,
	}
}

const trace = "BrandHandler"

func (bh *BrandHandler) GetBrands(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".GetBrands")

	brands, err := bh.useCase.GetBrands()
	if err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.SERVER_ERROR, errors.SERVER_ERROR_DESCR)
		return ctx.JSON(http.StatusInternalServerError, newBasketError)
	}

	if brands == nil {
		brands = make([]models.Brand, 0)
	}

	logger.Trace(trace + " success GetBrands")
	return ctx.JSON(http.StatusOK, brands)
}

func (bh *BrandHandler) GetProductsByBrand(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + "GetProductsByBrand")

	idString := ctx.QueryParam("id")
	if idString == "" {
		logger.Error("bad query param for GetProductsByBrand")
		err:= errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Error(err)
		err := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	products, err := bh.useCase.GetProductsByBrand(id)
	if err != nil {
		logger.Error(err)
		err := errors.NewError(errors.DB_ERROR, errors.SERVER_ERROR_DESCR)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	for i, _ := range products {
		products[i] = sanitizer.SanitizeData(&products[i]).(models.Product)
	}

	return ctx.JSON(http.StatusOK, products)
}
