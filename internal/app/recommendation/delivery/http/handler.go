package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation"
	sessionJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

const trace = "RecommendationHandler"

type RecommendHandler struct {
	useCase        recommendation.UseCase
	SessionManager sessionJwt.TokenManager
}

func NewRecommendHandler(useCase recommendation.UseCase, sessionManager sessionJwt.TokenManager) *RecommendHandler {
	return &RecommendHandler{
		useCase:        useCase,
		SessionManager: sessionManager,
	}
}

func (rh *RecommendHandler) GetRecommendation(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + "AddFavouriteProduct")

	var recommendProduct []models.Product
	idNum, err := rh.SessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err == nil {
		recommendProduct, err = rh.useCase.GetRecommendForUser(int(idNum))
		if err != nil {
			logger.Error(err)
			return ctx.JSON(http.StatusUnauthorized,
				errors.NewError(errors.DB_ERROR, errors.BD_ERROR_DESCR))
		}

		if len(recommendProduct) < 4 {
			diff := 4 - len(recommendProduct)
			tempProducts, err := rh.useCase.GetMostPopularProduct()
			tempProducts = tempProducts[:diff]
			for i := 0; i < diff; i++ {
				recommendProduct = append(recommendProduct, tempProducts[i])
			}
			if err != nil {
				logger.Error(err)
				return ctx.JSON(http.StatusBadRequest,
					errors.NewError(errors.DB_ERROR, errors.BD_ERROR_DESCR))
			}

		}
	} else {
		recommendProduct, err = rh.useCase.GetMostPopularProduct()
		if err != nil {
			logger.Error(err)
			return ctx.JSON(http.StatusUnauthorized,
				errors.NewError(errors.DB_ERROR, errors.BD_ERROR_DESCR))
		}
	}

	if recommendProduct == nil {
		recommendProduct = make([]models.Product, 0)
	}

	return ctx.JSON(http.StatusOK, recommendProduct)
}
