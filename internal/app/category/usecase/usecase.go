package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/usecase/helpers"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"strings"
)

type UseCase struct {
	repositoryCategory category.Repository
	repositoryProduct  product.Repository
}

func NewCategoryUseCase(repositoryCategory category.Repository, repositoryModel product.Repository) *UseCase {
	return &UseCase{
		repositoryCategory: repositoryCategory,
		repositoryProduct:  repositoryModel,
	}
}

func (uc *UseCase) GetProductsByCategory(filter postgre.Filter) (*models.ProductsCategory, error) {
	var result models.ProductsCategory

	products, err := uc.repositoryProduct.GetByCategory(filter)
	if err != nil {
		return nil, err
	}

	for i, _ := range products {
		imageSlice := strings.Split(products[i].Image, ";")
		products[i].Image = imageSlice[0]
	}

	result.Products = products

	minPrice, maxPrice, err := uc.repositoryCategory.GetMinMaxPriceCategory(filter.NameCategory)
	if err != nil {
		return nil, err
	}
	result.MinPrice = minPrice
	result.MaxPrice = maxPrice

	return &result, nil
}

func (uc *UseCase) GetAllCategories() (models.CategoryNode, error) {
	nestingCategories, err := uc.repositoryCategory.SelectAllCategories()
	if err != nil {
		return models.CategoryNode{}, err
	}

	node := helpers.ParseCategories(nestingCategories)
	return node, nil
}

func (uc *UseCase) CreateCategory(category models.CreateCategory) error {
	return uc.repositoryCategory.CreateCategory(category)
}
