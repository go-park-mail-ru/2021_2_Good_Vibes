package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"math"
)

type UseCase struct {
	repositorySearch search.Repository
}

func NewSearchUseCase(repositorySearch search.Repository) *UseCase {
	return &UseCase{
		repositorySearch: repositorySearch,
	}
}

func (uc *UseCase) GetSuggests(str string) (models.Suggest, error) {
	suggests, err := uc.repositorySearch.GetSuggests(str)
	if err != nil {
		return models.Suggest{}, err
	}

	return suggests, nil
}

func (uc *UseCase) GetSearchResults(searchString []string, filter postgre.Filter) (models.ProductsCategory, error) {
	var products []models.Product

	productMap := make(map[models.Product]int)

	resultProducts, err := uc.repositorySearch.GetSearchResults(searchString, filter)

	if err != nil {
		return models.ProductsCategory{}, err
	}

	for _, products := range resultProducts {
		for _, product := range products {
			productMap[product] += 1
		}
	}

	var maxValue int

	for _, value := range productMap {
		maxValue = int(math.Max(float64(maxValue), float64(value)))
	}
	fmt.Println(maxValue)

	for i := 0; i < maxValue; i++ {
		for key, value := range productMap {
			if value == maxValue-i {
				products = append(products, key)
			}
		}
	}

	var maxPrice int
	for _, product := range products {
		maxPrice = int(math.Max(float64(maxPrice), product.Price))
	}

	var minPrice int
	for _, product := range products {
		minPrice = int(math.Min(float64(minPrice), product.Price))
	}

	if minPrice == 0 {
		minPrice = maxPrice
	}

	productWithPrices := models.ProductsCategory{
		Products: products,
		MinPrice: float64(minPrice),
		MaxPrice: float64(maxPrice),
	}

	return productWithPrices, nil
}
