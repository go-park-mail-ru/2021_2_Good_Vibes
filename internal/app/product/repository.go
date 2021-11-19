package product

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
)

type Repository interface {
	Insert(prod models.Product) (int, error)
	GetAll() ([]models.Product, error)
	GetFavouriteProducts(userId int) ([]models.Product, error)
	GetById(id int) (models.Product, error)
	GetByCategory(filter postgre.Filter) ([]models.Product, error)
	SaveProductImageName(productId int, fileName string) error
	AddFavouriteProduct(product models.FavouriteProduct) error
	DeleteFavouriteProduct(product models.FavouriteProduct) error
}
