package panic

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				context.NoContent(http.StatusInternalServerError)
			}
		}()
		return next(context)
	}
}
