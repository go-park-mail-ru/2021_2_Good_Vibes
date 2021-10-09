package http

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/labstack/echo/v4"
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

func (ch *CategoryHandler) GetAllCategories(ctx echo.Context) error {
	val := ctx.QueryParams()

	nameString := val.Get("name")
	if nameString == "" {
		if AllCategoriesJson.Name != "" {
			return ctx.JSON(http.StatusOK, AllCategoriesJson)
		}

		categories, err := ch.useCase.GetAllCategories()
		if err != nil {
			newCategoryError := errors.NewError(errors.DB_ERROR, err.Error())
			return ctx.JSON(http.StatusBadRequest, newCategoryError)
		}
		AllCategoriesJson = categories

		return ctx.JSON(http.StatusOK, categories)
	}

	products, err := ch.useCase.GetProductsByCategory(nameString)
	if err != nil {
		newCategoryError := errors.NewError(errors.SERVER_ERROR, err.Error())
		return ctx.JSON(http.StatusBadRequest, newCategoryError)
	}

	return ctx.JSON(http.StatusOK, products)
}


