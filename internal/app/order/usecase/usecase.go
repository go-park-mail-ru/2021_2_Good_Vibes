package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order"
)

type UseCase struct {
	repositoryOrder order.Repository
}

func NewOrderUseCase(repositoryOrder order.Repository) *UseCase {
	return &UseCase{
		repositoryOrder: repositoryOrder,
	}
}

func (uc *UseCase) PutOrder(order models.Order) (int, error) {
	orderId, err := uc.repositoryOrder.PutOrder(order)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}
