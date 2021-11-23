package review

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	AddReview(review models.Review, productRating float64) error
	UpdateReview(review models.Review, productRating float64) error
	GetReviewsByProductId(productId int) ([]models.Review, error)
	GetAllRatingsOfProduct(productId int) ([]models.ProductRating, error)
	GetReviewsByUser(userName string) ([]models.Review, error)
	GetReviewByUserAndProduct(userId int, productId int) (models.Review, error)
	DeleteReview(userId int, productId int, productRating float64) error
}
