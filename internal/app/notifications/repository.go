package notifications

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type Repository interface {
	GetStatusChanges() ([]models.ChangedStatus, error)
	GetAddressInfo(orderId int) (models.Address, error)
}
