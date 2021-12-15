package http

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review"
	sessionJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt"
	customLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ReviewHandler struct {
	useCase        review.UseCase
	sessionManager sessionJwt.TokenManager
}

func NewReviewHandler(useCase review.UseCase, sessionManager sessionJwt.TokenManager) *ReviewHandler {
	return &ReviewHandler{
		useCase:        useCase,
		sessionManager: sessionManager,
	}
}

const trace = "ReviewHandler"

func (rh *ReviewHandler) AddReview(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".AddReview")

	var newReview models.Review
	userId, err := rh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))
	}

	newReview.UserId = int(userId)

	if err := ctx.Bind(&newReview); err != nil {
		logger.Error(err)
		newError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err := ctx.Validate(&newReview); err != nil {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	err = rh.useCase.AddReview(&newReview)
	if err != nil && err.Error() == customErrors.REVIEW_EXISTS_DESCR {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.REVIEW_EXISTS_ERROR, customErrors.REVIEW_EXISTS_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err != nil {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newError)
	}

	logger.Trace(trace + " success ReviewHandler")

	return ctx.JSON(http.StatusOK, newReview)
}

func (rh *ReviewHandler) UpdateReview(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".UpdateReview")

	var newReview models.Review
	userId, err := rh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))
	}

	newReview.UserId = int(userId)

	if err := ctx.Bind(&newReview); err != nil {
		logger.Error(err)
		newError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err := ctx.Validate(&newReview); err != nil {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	err = rh.useCase.UpdateReview(newReview)
	if err != nil && err.Error() == customErrors.NO_REVIEW_DESCR {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.NO_REVIEW_ERROR, customErrors.NO_REVIEW_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err != nil {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newError)
	}

	logger.Trace(trace + " success ReviewHandler")

	return ctx.JSON(http.StatusOK, newReview)
}

func (rh *ReviewHandler) DeleteReview(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " DeleteReview")

	userId, err := rh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))
	}

	var productId models.ProductId

	if err := ctx.Bind(&productId); err != nil {
		logger.Error(err)
		newError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err := ctx.Validate(&productId); err != nil {
		logger.Error(err, productId)
		newError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	err = rh.useCase.DeleteReview(int(userId), productId.ProductId)
	if err != nil && err.Error() == customErrors.NO_REVIEW_DESCR {
		logger.Error(err, productId)
		newError := customErrors.NewError(customErrors.NO_REVIEW_ERROR, customErrors.NO_REVIEW_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	logger.Trace(trace + " success DeleteReview")

	return ctx.JSON(http.StatusOK, productId)
}

func (rh *ReviewHandler) GetReviewsByProductId(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetReviewsByProductId")

	productIdString := ctx.QueryParam("product_id")

	if productIdString == "" {
		logger.Error("bad query param for GetReviewsByProductId")
		newError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	productId, err := strconv.Atoi(productIdString)

	reviews, err := rh.useCase.GetReviewsByProductId(productId)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if reviews == nil {
		reviews = make([]models.Review, 0)
	}

	logger.Trace(trace + " success GetReviewsByProductId")
	return ctx.JSON(http.StatusOK, reviews)
}

func (rh *ReviewHandler) GetReviewsByUser(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetReviewsByUser")

	userName := ctx.QueryParam("name")

	if userName == "" {
		logger.Error("bad query param for GetReviewsByUser")
		newError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newError)
	}

	reviews, err := rh.useCase.GetReviewsByUser(userName)

	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	if reviews == nil {
		reviews = make([]models.Review, 0)
	}

	logger.Trace(trace + " success GetReviewsByUser")
	return ctx.JSON(http.StatusOK, reviews)
}
