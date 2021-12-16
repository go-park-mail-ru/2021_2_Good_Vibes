package recommendation

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	GetRecommendForUser(userId int) ([]models.Product, error)
	GetMostPopularProduct()([]models.Product, error)
}
