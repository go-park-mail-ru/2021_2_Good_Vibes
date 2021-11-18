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

func (rh *ReviewHandler) AddReview (ctx echo.Context) error {
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

	err = rh.useCase.AddReview(newReview)
	if err != nil && err.Error() == customErrors.RATING_EXISTS_DESCR {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.RATING_EXISTS_ERROR, customErrors.RATING_EXISTS_DESCR)
		return ctx.JSON(http.StatusOK, newError)
	}

	if err != nil {
		logger.Error(err, newReview)
		newError := customErrors.NewError(customErrors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusInternalServerError, newError)
	}

	logger.Trace(trace + " success ReviewHandler")

	return ctx.JSON(http.StatusOK, newReview)
}


func (rh *ReviewHandler) GetReviewsByProductId(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetReviewsByProductId")

	productIdString:= ctx.QueryParam("product_id")

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

	logger.Trace(trace + " success GetReviewsByProductId")
	return ctx.JSON(http.StatusOK, reviews)
}


func (rh *ReviewHandler) GetReviewsByUser(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + " GetReviewsByUser")

	userName:= ctx.QueryParam("name")

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

	logger.Trace(trace + " success GetReviewsByUser")
	return ctx.JSON(http.StatusOK, reviews)
}

/*
func (bh *BasketHandler) PutInBasket(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".PutInBasket")

	var newProduct models.BasketProduct
	userId, err := bh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	newProduct.UserId = int(userId)

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
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".GetBasket")

	var user models.UserID

	userId, err := bh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	user.UserId = int(userId)

	var basketProducts []models.BasketProduct

	basketProducts, err = bh.useCase.GetBasket(user.UserId)
	if err != nil {
		logger.Error(err)
		newBasketError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newBasketError)
	}

	if basketProducts == nil {
		basketProducts = make([]models.BasketProduct, 0)
	}

	logger.Trace(trace + " success GetBasket")
	return ctx.JSON(http.StatusOK, basketProducts)
}

func (bh *BasketHandler) DropBasket(ctx echo.Context) error {
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".DropBasket")

	var user models.UserID

	userId, err := bh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	user.UserId = int(userId)

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
	logger := customLogger.TryGetLoggerFromContext(ctx)
	logger.Trace(trace + ".DeleteProduct")

	var product models.BasketProduct

	userId, err := bh.sessionManager.ParseTokenFromContext(ctx.Request().Context())
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnauthorized, errors.NewError(errors.TOKEN_ERROR, errors.TOKEN_ERROR_DESCR))
	}

	product.UserId = int(userId)

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
*/