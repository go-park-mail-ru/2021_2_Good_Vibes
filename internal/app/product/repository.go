package product

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	Insert(prod models.Product) (int, error)
	PutSalesProduct(sales models.SalesProduct) error
	GetAll() ([]models.Product, error)
	GetNewProducts() ([]models.Product, error)
	GetSalesProducts() ([]models.Product, error)
	GetProductsByBrand(brandId int) ([]models.Product, error)
	GetFavouriteProducts(userId int) ([]models.Product, error)
	GetById(id int) (models.Product, error)
	GetByCategory(filter postgre.Filter) ([]models.Product, error)
	SaveProductImageName(productId int, fileName string) error
	AddFavouriteProduct(product models.FavouriteProduct) error
	DeleteFavouriteProduct(product models.FavouriteProduct) error
	ChangeRecommendUser(userId int, ProductId int, isSearch string) error
	TryGetProductWithSimilarName(productName string) ([]models.Product, error)
	IsFavourite(productID int, userID int64) (*bool, error)
}
