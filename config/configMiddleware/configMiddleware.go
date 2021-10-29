package configMiddleware

import (
	middlewareCors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/cors"
	middlewareCSRF "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/csrf"
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
		middlewareCSRF.CSRFWithConfig(),
	)
}
