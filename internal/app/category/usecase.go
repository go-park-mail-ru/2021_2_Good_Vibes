package category

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
)

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	GetAllCategories() (models.CategoryNode, error)
	GetProductsByCategory(filter postgre.Filter) ([]models.Product, error)
	CreateCategory(category models.CreateCategory) error
}
