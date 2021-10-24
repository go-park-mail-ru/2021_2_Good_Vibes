package manager

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"time"
)
type Manager struct {
	secretKey string
}

const (
	FieldNameId = "id"
	FieldNameTime = "exp"
)

func NewTokenManager (secretKey string) (*Manager, error) {
	if secretKey == "" {
		return nil, errors.New(customErrors.BAD_INIT_SECRET_KEY)
	}
	return &Manager{secretKey: secretKey}, nil
}

func (m *Manager)GetToken(id int, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		FieldNameId: id,
		FieldNameTime : time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) ParseToken(accessToken string) (string, error) {
	return "sdfs", nil
}