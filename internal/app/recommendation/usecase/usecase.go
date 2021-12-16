package usecase

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation"
	"math/rand"
)

type RecommendationUseCase struct {
	repository        recommendation.Repository
	repositoryProduct product.Repository
}

func NewRecommendationUseCase(repository recommendation.Repository, repositoryProduct product.Repository) *RecommendationUseCase {
	return &RecommendationUseCase{
		repository:        repository,
		repositoryProduct: repositoryProduct,
	}
}

func (ru *RecommendationUseCase) GetRecommendForUser(userId int) ([]models.Product, error) {
	var recommendProductModels []models.Product
	recommendProductsAll, err := ru.repository.GetRecommendProductForUser(userId)

	var randomArray []int
	if len(recommendProductsAll) < 4 {
		for index, _ := range recommendProductsAll {
			randomArray = append(randomArray, index)
		}
	} else {
		randomArray = rand.Perm(len(recommendProductsAll))[:4]
	}

	for _, value := range randomArray {
		currentId := recommendProductsAll[value].Id
		currentProduct, err := ru.repositoryProduct.GetById(currentId)
		if err != nil {
			return nil, err
		}
		recommendProductModels = append(recommendProductModels, currentProduct)
	}

	if err != nil {
		return nil, err
	}
	return recommendProductModels, nil
}

func (ru *RecommendationUseCase) GetMostPopularProduct() ([]models.Product, error) {
	var fourMostPopularProduct []models.Product
	MostPopularProducts, err := ru.repository.GetMostPopularProduct()
	if err != nil {
		return nil, err
	}

	randomArray := rand.Perm(len(MostPopularProducts))[:4]

	for _, value := range randomArray {
		fourMostPopularProduct = append(fourMostPopularProduct, MostPopularProducts[value])
	}

	return fourMostPopularProduct, nil
}
