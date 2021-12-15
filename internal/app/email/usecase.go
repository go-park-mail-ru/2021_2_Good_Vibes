package email

type UseCase interface {
	ConfirmEmail(email string, token string) error
	SendConfirmationEmail(email string, token string) error
	GetUserEmailById(id uint64) (string, error)
	GenerateToken() string
	InsertUserToken(email string, token string) error
}
