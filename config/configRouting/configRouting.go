package configRouting

import (
	basketHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	middlewareAut "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/authentication"
	orderHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	handler2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	http2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServerConfigRouting struct {
	UserHandler    *http2.UserHandler
	ProductHandler *handler2.ProductHandler
	OrderHandler *orderHttp.OrderHandler
	BasketHandler *basketHttp.BasketHandler
}

func (cr *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.Static("/", "static")
	router.POST("/login", cr.UserHandler.Login)
	router.POST("/signup", cr.UserHandler.SignUp)
	router.GET("/profile", profile, middlewareAut.IsLogin)
	router.GET("/homepage", cr.ProductHandler.GetAllProducts)
	router.GET("/product", cr.ProductHandler.GetProductById)
	router.GET("/logout", cr.UserHandler.Logout, middlewareAut.IsLogin)
	router.POST("/order/confirm", cr.OrderHandler.PutOrder)
	router.POST("/order/basket/put", cr.BasketHandler.PutInBasket)
	router.POST("/order/basket/drop", cr.BasketHandler.DropBasket)
	router.POST("/order/basket/delete", cr.BasketHandler.DeleteProduct)
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
