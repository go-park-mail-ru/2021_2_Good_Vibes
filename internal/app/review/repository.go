package review

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	AddReview(review models.Review, productRating float64) error
	GetReviewsByProductId(productId int) ([]models.Review, error)
	GetAllRatingsOfProduct(productId int) ([]models.ProductRating, error)
	GetReviewsByUser(userName string) ([]models.Review, error)
}
