package cors

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var (
	allowOrigins  = []string{"https://dreamy-yonath-26f2eb.netlify.app", "https://peaceful-bell-d76220.netlify.app"}
	allowMethods  = []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions}
	allowHeaders  = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Csrf-Token"}
	exposeHeaders = []string{"Authorization", "X-Csrf-Token"}
)

func GetCORSConfigStruct() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: true,
	}
}
