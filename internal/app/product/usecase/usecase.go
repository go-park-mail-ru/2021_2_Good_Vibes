package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	guuid "github.com/google/uuid"
)

type UseCase struct {
	repository product.Repository
}

func NewProductUsecase(repositoryProduct product.Repository) *UseCase {
	return &UseCase{
		repository: repositoryProduct,
	}
}

func (uc *UseCase) AddProduct(product models.Product) (int, error) {
	return uc.repository.Insert(product)
}

func (uc *UseCase) GetAllProducts() ([]models.Product, error) {
	return uc.repository.GetAll()
}

func (uc *UseCase) AddFavouriteProduct(product models.FavouriteProduct) error {
	return uc.repository.AddFavouriteProduct(product)
}

func (uc *UseCase) DeleteFavouriteProduct(product models.FavouriteProduct) error {
	return uc.repository.DeleteFavouriteProduct(product)
}

func (uc *UseCase) GetFavouriteProducts(userId int) ([]models.Product, error) {
	return uc.repository.GetFavouriteProducts(userId)
}

func (uc *UseCase) GetProductById(id int) (models.Product, error) {
	return uc.repository.GetById(id)
}

func (uc *UseCase) GenerateProductImageName() string {
	return guuid.New().String()
}

func (uc *UseCase) SaveProductImageName(productId int, fileName string) error {
	err := uc.repository.SaveProductImageName(productId, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) ChangeRecommendUser(userId int, ProductId int, isSearch string) error {
	err := uc.repository.ChangeRecommendUser(userId, ProductId, isSearch)
	if err != nil {
		return err
	}
	return nil
}
