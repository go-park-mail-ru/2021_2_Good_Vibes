package brands

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	GetBrands() ([]models.Brand, error)
	GetProductsByBrand(id int) (models.ProductsBrand, error)
}
