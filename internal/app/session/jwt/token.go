package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"strconv"
	"time"
)

func GetToken(id int, name string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = strconv.Itoa(id)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	return token.SignedString([]byte(config.ConfigApp.MainConfig.SecretKey))
}
