package storage

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"

type UseCase interface {
	AddProduct(prod product.Product) (int, error)
	// это наверно в будущем не нужно
	GetAllProducts() ([]product.Product, error)
	// GetProductsOnPage(page int) ([]product.Product, error)
	GetProductById(id int) (product.Product, error)
}
