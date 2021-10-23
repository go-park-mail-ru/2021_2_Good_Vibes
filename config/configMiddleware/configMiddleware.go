package configMiddleware

import (
	middlewareCors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/cors"
	middlewarePanic "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/panic"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigMiddleware(router *echo.Echo)  {
	router.Use(
		middlewarePanic.Recover,
		middleware.CORSWithConfig(middlewareCors.GetCORSConfigStruct()),
		)
}