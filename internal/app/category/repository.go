package category

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	SelectAllCategories() ([]models.NestingCategory, error)
	CreateCategory(category models.CreateCategory) error
	GetMinMaxPriceCategory(category string) (float64,float64,error)
}
