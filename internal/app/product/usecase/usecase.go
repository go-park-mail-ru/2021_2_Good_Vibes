package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
)

type Usecase struct {
	repository product.Repository
}

func NewProductUsecase(repositoryProduct product.Repository) *Usecase {
	return &Usecase{
		repository: repositoryProduct,
	}
}

func (uc *Usecase) AddProduct(prod models.Product) error {
	return uc.repository.Insert(prod)
}

func (uc *Usecase) GetAllProducts() ([]models.Product, error) {
	return uc.repository.GetAll()
}

func (uc *Usecase) GetProductById(id int) (models.Product, error) {
	return uc.repository.GetById(id)
}
