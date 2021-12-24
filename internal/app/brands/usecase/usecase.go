package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"strings"
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

func (uc *UseCase) GetProductsByBrand(id int) (models.ProductsBrand, error) {
	var result models.ProductsBrand
	products, err := uc.repositoryProduct.GetProductsByBrand(id)
	if len(products) != 0 {
		result.BrandName = products[0].BrandName
	}
	if err != nil {
		return models.ProductsBrand{}, err
	}

	for i, _ := range products {
		imageSlice := strings.Split(products[i].Image, ";")
		products[i].Image = imageSlice[0]
	}
	result.Products = products
	return result, err
}
