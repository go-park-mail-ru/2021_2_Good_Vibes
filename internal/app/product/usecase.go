package product

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Usecase interface {
	AddProduct(prod models.Product) error
	// это наверно в будущем не нужно
	GetAllProducts() ([]models.Product, error)
	// GetProductsOnPage(page int) ([]product.Product, error)
	GetProductById(id int) (models.Product, error)
}
