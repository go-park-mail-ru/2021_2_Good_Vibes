package category

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	GetAllCategories() (models.CategoryNode, error)
	GetProductsByCategory(category string) ([]models.Product, error)
	CreateCategory(categoryName string, parentCategoryName string) error
}
