package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
)

type UseCase struct {
	repositoryBrand brands.Repository
	repositoryProduct product.Repository
}

func NewBrandUseCase(repositoryBrand brands.Repository, repositoryProduct product.Repository) *UseCase {
	return &UseCase{
		repositoryBrand: repositoryBrand,
		repositoryProduct: repositoryProduct,
	}
}

func (uc *UseCase) GetBrands() ([]models.Brand, error) {
	return uc.repositoryBrand.GetBrands()
}

func (uc *UseCase) GetProductsByBrand(id int) ([]models.Product, error) {
	return uc.repositoryProduct.GetProductsByBrand(id)
}
