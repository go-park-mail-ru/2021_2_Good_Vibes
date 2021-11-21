package order

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	PutOrder(order models.Order) (int, float64, error)
	GetAllOrders(user int) ([]models.Order, error)
}
