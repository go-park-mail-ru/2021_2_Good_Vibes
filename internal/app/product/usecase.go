package product

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	AddProduct(prod models.Product) (int, error)
	GetAllProducts() ([]models.Product, error)
	// GetProductsOnPage(page int) ([]product.Product, error)
	GetProductById(id int) (models.Product, error)
	GenerateProductImageName() string
	SaveProductImageName(productId int, fileName string) error
}
