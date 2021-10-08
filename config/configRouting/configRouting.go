package configRouting

import (
	middlewareAut "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/authentication"
	handler2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/handler"
	http2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServerConfigRouting struct {
	UserHandler    *http2.UserHandler
	ProductHandler *handler2.ProductHandler
}

func (cr *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.Static("/", "static")
	router.POST("/login", cr.UserHandler.Login)
	router.POST("/signup", cr.UserHandler.SignUp)
	router.GET("/profile", profile, middlewareAut.IsLogin)
	router.GET("/homepage", cr.ProductHandler.GetAllProducts)
	router.GET("/product", cr.ProductHandler.GetProductById)
	router.GET("/logout", cr.UserHandler.Logout, middlewareAut.IsLogin)
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
