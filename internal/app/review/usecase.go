package review

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	AddReview(review models.Review) error
}
