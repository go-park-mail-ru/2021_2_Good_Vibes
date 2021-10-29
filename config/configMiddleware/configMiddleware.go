package configMiddleware

import (
	middlewareCors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/cors"
	middlewareLogger "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/logger"
	middlewarePanic "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/panic"
	middlewareRequestId "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/requestId"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigMiddleware(router *echo.Echo) {
	router.Use(
		middlewarePanic.Recover,
		middlewareRequestId.RequestId,
		middlewareLogger.AccessLog,
		middleware.CORSWithConfig(middlewareCors.GetCORSConfigStruct()),
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper: func(context echo.Context) bool {
				if context.Request().RequestURI == "/login" ||
					context.Request().RequestURI == "/signup" {
					return true
				}
				return false
			},
			CookieHTTPOnly: true,
		}),
		middleware.Secure(),
	)
}
