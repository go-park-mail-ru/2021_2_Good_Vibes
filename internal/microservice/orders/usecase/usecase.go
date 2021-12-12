package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders"
)

type UseCase struct {
	repositoryOrder orders.Repository
}

func NewOrderUseCase(repositoryOrder orders.Repository) *UseCase {
	return &UseCase{
		repositoryOrder: repositoryOrder,
	}
}

func (uc *UseCase) PutOrder(order models.Order) (int, float64, error) {
	productPrices, err := uc.repositoryOrder.SelectPrices(order.Products)

	if err != nil {
		return 0, 0, err
	}

	productPricesMap := make(map[int]float64, len(productPrices))
	for _, productPrice := range productPrices {
		productPricesMap[productPrice.Id] = productPrice.Price
	}

	changePriceAfterParcing := 0
	if order.Promocode != "" {
		promoCode, err := uc.repositoryOrder.CheckPromoCode(order.Promocode)
		if promoCode != nil && err == nil {
			changePriceAfterParcing = uc.ParsePromoCode(*promoCode, productPricesMap)
		}
	}
	fmt.Println("mapa after promo:", productPricesMap)

	var cost float64
	for index, product := range order.Products {
		TotalPriceProduct := float64(product.Number) * productPricesMap[product.ProductId]
		order.Products[index].Price = TotalPriceProduct
		cost += TotalPriceProduct
	}

	fmt.Println("cost before promo", cost)
	cost -= float64(changePriceAfterParcing)
	if cost < 1 {
		cost = 1
	}
	order.Cost = cost
	fmt.Println("cost after promo", cost)

	orderId, err := uc.repositoryOrder.PutOrder(order)
	if err != nil {
		return 0, 0, err
	}

	return orderId, cost, nil
}

func (uc *UseCase) CheckOrder(order models.Order) (*models.Order, error) {
	productPrices, err := uc.repositoryOrder.SelectPrices(order.Products)
	if err != nil {
		return nil, err
	}

	productPricesMap := make(map[int]float64, len(productPrices))
	for index, productPrice := range productPrices {
		productPricesMap[productPrice.Id] = productPrice.Price
		order.Products[index].Price = productPrice.Price
	}

	changePriceAfterParcing := 0
	if order.Promocode != "" {
		promoCode, err := uc.repositoryOrder.CheckPromoCode(order.Promocode)
		if promoCode != nil && err == nil {
			changePriceAfterParcing = uc.ParsePromoCode(*promoCode, productPricesMap)
		}
	}
	fmt.Println("mapa after promo:", productPricesMap)

	var cost float64
	for index, product := range order.Products {
		TotalPriceProduct := float64(product.Number) * productPricesMap[product.ProductId]
		order.Products[index].PriceWithPromo = TotalPriceProduct
		cost += TotalPriceProduct
	}

	fmt.Println("cost before promo", cost)
	cost -= float64(changePriceAfterParcing)
	if cost < 1 {
		cost = 1
	}
	order.Cost = cost
	fmt.Println("cost after promo", cost)

	return &order, nil
}

func (uc *UseCase) GetAllOrders(user int) ([]models.Order, error) {
	return uc.repositoryOrder.GetAllOrders(user)
}

func (uc *UseCase) ParsePromoCode(code models.PromoCode, prices map[int]float64) int {
	if code.CategoryId == -1 && code.ProductId == -1 &&
		code.Type == models.TypePromoInterest {
		for key, value := range prices {
			prices[key] = value - value/100*float64(code.Value)
		}
		return 0
	}

	if code.CategoryId == -1 && code.ProductId == -1 &&
		code.Type == models.TypePromoMoney {
		return code.Value
	}

	if code.CategoryId == -1 && code.ProductId != -1 &&
		code.Type == models.TypePromoInterest {
		if value, ok := prices[code.ProductId]; ok {
			prices[code.ProductId] = value - value/100*float64(code.Value)
		} else {
			return -1
		}
		return 0
	}

	if code.CategoryId == -1 && code.ProductId != -1 &&
		code.Type == models.TypePromoMoney {
		if _, ok := prices[code.ProductId]; ok {
			prices[code.ProductId] -= float64(code.Value)
		} else {
			return -1
		}

		if prices[code.ProductId] < 1 {
			prices[code.ProductId] = 1
		}
		return 0
	}

	if code.CategoryId != -1 && code.ProductId == -1 &&
		code.Type == models.TypePromoInterest {
		for key, value := range prices {
			categoryProduct, err := uc.repositoryOrder.GetProductCategory(key)
			fmt.Println(err)
			if categoryProduct == code.CategoryId {
				prices[key] = value - value/100*float64(code.Value)
			}
		}
		return 0
	}

	if code.CategoryId != -1 && code.ProductId == -1 &&
		code.Type == models.TypePromoMoney {
		for key, _ := range prices {
			categoryProduct, err := uc.repositoryOrder.GetProductCategory(key)
			fmt.Println(err)
			if categoryProduct == code.CategoryId {
				prices[key] -= float64(code.Value)
			}
			if prices[key] < 1 {
				prices[key] = 1
			}
		}
		return 0
	}

	return 0
}
