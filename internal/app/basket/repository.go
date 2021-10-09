package basket

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	PutInBasket(basket models.BasketProduct) error
	DropBasket(userId int) error
	DeleteProduct(product models.BasketProduct) error
}
