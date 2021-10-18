package configRouting

import (
	basketHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	categoryHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/delivery/http"
	middlewareAut "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/authentication"
	orderHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	handler2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	http2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ServerConfigRouting struct {
	UserHandler     *http2.UserHandler
	ProductHandler  *handler2.ProductHandler
	OrderHandler    *orderHttp.OrderHandler
	BasketHandler   *basketHttp.BasketHandler
	CategoryHandler *categoryHttp.CategoryHandler
}

func (cr *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.Static("/", "static")
	router.POST("/login", cr.UserHandler.Login)
	router.POST("/signup", cr.UserHandler.SignUp)
	router.GET("/myprofile", cr.UserHandler.Profile, middlewareAut.IsLogin)
	router.GET("/homepage", cr.ProductHandler.GetAllProducts)
	router.GET("/product", cr.ProductHandler.GetProductById)
	router.GET("/logout", cr.UserHandler.Logout, middlewareAut.IsLogin)
	router.GET("/cart/get", cr.BasketHandler.GetBasket, middlewareAut.IsLogin)
	router.POST("/cart/confirm", cr.OrderHandler.PutOrder, middlewareAut.IsLogin)
	router.POST("/cart/put", cr.BasketHandler.PutInBasket, middlewareAut.IsLogin)
	router.POST("/cart/drop", cr.BasketHandler.DropBasket, middlewareAut.IsLogin)
	router.POST("/cart/delete", cr.BasketHandler.DeleteProduct, middlewareAut.IsLogin)
	router.GET("/category", cr.CategoryHandler.GetCategories)
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
