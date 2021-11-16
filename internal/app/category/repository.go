package category

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	SelectAllCategories() ([]models.NestingCategory, error)
	CreateCategory(category models.CreateCategory) error
}
