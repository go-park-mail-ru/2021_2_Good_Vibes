package basket

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	PutInBasket(basket models.BasketProduct) error
	GetBasket(userId int) ([]models.BasketProduct, error)
	DropBasket(userId int) error
	DeleteProduct(product models.BasketProduct) error
}
