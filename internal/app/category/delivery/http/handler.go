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

func (ch *CategoryHandler) GetCategories(ctx echo.Context) error {
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

func (ch *CategoryHandler) CreateCategory(ctx echo.Context) error {
	var newCategory models.CreateCategory

	if err := ctx.Bind(&newCategory); err != nil {
		newSignupError := errors.NewError(errors.BIND_ERROR, errors.BIND_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	if err := ctx.Validate(&newCategory); err != nil {
		newSignupError := errors.NewError(errors.VALIDATION_ERROR, errors.VALIDATION_DESCR)
		return ctx.JSON(http.StatusBadRequest, newSignupError)
	}

	err := ch.useCase.CreateCategory(newCategory.Category, newCategory.ParentCategory)
	if err != nil {
		return err
	}

	categories, err := ch.useCase.GetAllCategories()
	if err != nil {
		return err
	}

	AllCategoriesJson = categories

	return ctx.JSON(http.StatusOK, categories)
}
