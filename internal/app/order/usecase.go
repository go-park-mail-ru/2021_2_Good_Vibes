package order

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	PutOrder(order models.Order) (int, float64, error)
	GetAllOrders(user int) ([]models.Order, error)
	GetOrderPriceWithPromo(order models.Order)(*models.Order, error)
}
