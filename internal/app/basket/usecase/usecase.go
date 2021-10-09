package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type UseCase struct {
	repositoryBasket basket.Repository
}

func NewBasketUseCase(repositoryBasket basket.Repository) *UseCase {
	return &UseCase{
		repositoryBasket: repositoryBasket,
	}
}

func (uc *UseCase) PutInBasket(basket models.BasketProduct) error {
	err := uc.repositoryBasket.PutInBasket(basket)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DropBasket(userId int) error {
	err := uc.repositoryBasket.DropBasket(userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DeleteProduct(product models.BasketProduct) error {
	err := uc.repositoryBasket.DeleteProduct(product)
	if err != nil {
		return err
	}

	return nil
}
