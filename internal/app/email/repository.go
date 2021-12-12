package email

type Repository interface {
	ConfirmEmail(email string) error
    GetToken(email string) (string, error)
}
