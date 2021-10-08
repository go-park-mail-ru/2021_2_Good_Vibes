package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func IsLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		cookie, err := context.Cookie("session_id")
		if err != nil {
			return context.NoContent(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signin method")
			}
			return []byte(config.ConfigApp.SecretKey), nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return next(context)
		}

		return context.NoContent(http.StatusUnauthorized)
	}
}
