package email

type UseCase interface {
	ConfirmEmail(email string, token string) error
}
