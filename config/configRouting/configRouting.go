package configRouting

import (
	handler2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/handler"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/handler"
	middleware_user "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/middleware"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServerConfigRouting struct {
	UserHandler *handler.UserHandler
	ProductHandler	*handler2.ProductHandler
}


func (cr *ServerConfigRouting) ConfigRouting (router *echo.Echo) {
	router.Static("/", "static")
	router.POST("/login", cr.UserHandler.Login)
	router.POST("/signup", cr.UserHandler.SignUp)
	router.GET("/profile", profile, middleware_user.IsLogin)
	router.GET("/homepage", cr.ProductHandler.GetAllProducts)
	router.GET("/product", cr.ProductHandler.GetProductById)
	router.GET("/logout", cr.UserHandler.Logout, middleware_user.IsLogin)
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}