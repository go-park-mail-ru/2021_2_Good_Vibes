package notifications

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	SearchStatusChanges() error
	ServeStatus(change models.ChangedStatus) error
	SendEmail(notifyInfo models.NotifyInfo) error
}
