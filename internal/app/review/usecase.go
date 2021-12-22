package review

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type UseCase interface {
	AddReview(review *models.Review) error
	UpdateReview(review models.Review) error
	GetReviewsByProductId(productId int) ([]models.Review, error)
	GetReviewsByUser(userName string) ([]models.Review, error)
	DeleteReview(userId int, productId int) error
}
