package handler

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/storage"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	storageProd storage.UseCase
}

func NewProductHandler(useCase *storage.UseCase) *ProductHandler {
	return &ProductHandler{
		storageProd: *useCase,
	}
}

//это должно быть, пока не нужно
func addProduct(ctx echo.Context) error {
	return nil
}

//пока решаем вопросы с пагинацией - так
func (ph *ProductHandler) GetAllProducts(ctx echo.Context) error {
	answer, err := ph.storageProd.GetAllProducts()
	if err != nil {
		newProductError := product.NewError(errors.DB_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}
	return ctx.JSON(http.StatusOK, answer)
}

func (ph *ProductHandler) GetProductById(ctx echo.Context) error {
	val := ctx.QueryParams()

	idString := val.Get("id")
	if idString == "" {
		newProductError := product.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newProductError)
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		newProductError := product.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadGateway, newProductError)
	}

	answer, err := ph.storageProd.GetProductById(id)

	if  err != nil {
		newProductError := product.NewError(errors.DB_ERROR, err.Error())
		return ctx.JSON(http.StatusOK, newProductError)
	}

	return ctx.JSON(http.StatusOK, answer)
}
