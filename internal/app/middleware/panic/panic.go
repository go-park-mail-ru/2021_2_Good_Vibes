package panic

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				reqId := fmt.Sprintf("%v", context.Get("reqId"))
				err_ := fmt.Sprintf("%v", err)
				logger.CustomLogger.LogErrorInfo(reqId, err_)

				context.NoContent(http.StatusInternalServerError)
			}
		}()

		return next(context)
	}
}
