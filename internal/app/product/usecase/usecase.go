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

func (uc *UseCase) PutSalesForProduct(sales models.SalesProduct) error {
	return uc.repository.PutSalesProduct(sales)
}

func (uc *UseCase) GetAllProducts() ([]models.Product, error) {
	return uc.repository.GetAll()
}

func (uc *UseCase) GetSalesProducts() ([]models.Product, error) {
	return uc.repository.GetSalesProducts()
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

func (uc *UseCase) GetProductById(id int, userID int64) (models.Product, error) {
	prod, err := uc.repository.GetById(id)
	if err != nil {
		return models.Product{}, err
	}

	if prod.Id == 0 {
		return models.Product{}, nil
	}

	if userID == 0 {
		return prod, nil
	}

	isFavourite, err := uc.repository.IsFavourite(id, userID)
	if err != nil {
		return models.Product{}, err
	}

	prod.IsFavourite = isFavourite

	return  prod, nil
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
