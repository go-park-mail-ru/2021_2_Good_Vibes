package jwt

type TokenManager interface {
	GetToken(id int, name string) (string, error)
	ParseToken(accessToken string) (string, error)
}

