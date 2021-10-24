package manager

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/labstack/echo/v4"
	"strconv"
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

func (m *Manager) ParseTokenFromContext(ctx echo.Context) (uint64, error) {
	token, ok := ctx.Get("token").(*jwt.Token)
	if !ok {
		return customErrors.TOKEN_ERROR, errors.New(customErrors.TOKEN_ERROR_DESCR)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return customErrors.TOKEN_ERROR, errors.New(customErrors.TOKEN_ERROR_DESCR)
	}

	idString := claims["id"].(string)
	idNum, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return customErrors.TOKEN_ERROR, err
	}

	return idNum, nil
}