package product

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	Insert(prod models.Product) (int, error)
	GetAll() ([]models.Product, error)
	GetById(id int) (models.Product, error)
	GetByCategory(category string) ([]models.Product, error)
	SaveProductImageName(productId int, fileName string) error
}
