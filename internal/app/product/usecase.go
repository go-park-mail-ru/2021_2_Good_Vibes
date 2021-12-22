package product

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	AddProduct(prod models.Product) (int, error)
	GetAllProducts() ([]models.Product, error)
	GetNewProducts() ([]models.Product, error)
	GetSalesProducts() ([]models.Product, error)
	PutSalesForProduct(sales models.SalesProduct) error
	// GetProductsOnPage(page int) ([]product.Product, error)
	GetFavouriteProducts(userId int) ([]models.Product, error)
	GetProductById(id int, userID int64) (models.Product, error)
	GenerateProductImageName() string
	SaveProductImageName(productId int, fileName string) error
	AddFavouriteProduct(product models.FavouriteProduct) error
	DeleteFavouriteProduct(product models.FavouriteProduct) error
	ChangeRecommendUser(userId int, ProductId int, isSearch string) error
}
