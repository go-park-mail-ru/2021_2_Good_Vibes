package handler

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/storage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProductHandler struct {
	storageProd storage.UseCase
}

func NewProductHandler(useCase storage.UseCase) *ProductHandler {
	return &ProductHandler{
		storageProd: useCase,
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
		ctx.NoContent(http.StatusBadRequest)
	}
	return ctx.JSON(http.StatusOK, answer)
}
