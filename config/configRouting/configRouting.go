package configRouting

import (
	basketHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	brandsHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands/delivery/http"
	categoryHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/delivery/http"
	middlewareAut "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/middleware/authentication"
	orderHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	handler2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	recommendation "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation/delivery/http"
	reviewHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/delivery/http"
	searchHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search/delivery/http"
	userHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerConfigRouting struct {
	UserHandler      *userHttp.UserHandler
	ProductHandler   *handler2.ProductHandler
	OrderHandler     *orderHttp.OrderHandler
	BasketHandler    *basketHttp.BasketHandler
	CategoryHandler  *categoryHttp.CategoryHandler
	ReviewHandler    *reviewHttp.ReviewHandler
	SearchHandler    *searchHttp.SearchHandler
	RecommendHandler *recommendation.RecommendHandler
	BrandHandler    *brandsHttp.BrandHandler
}

func (cr *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	// router.Static("/img/avatar", "upload/img/avatars")
	router.POST("/api/login", cr.UserHandler.Login)
	router.POST("/api/signup", cr.UserHandler.SignUp)
	router.POST("/api/upload/avatar", cr.UserHandler.UploadAvatar, middlewareAut.IsLogin)
	router.GET("/api/profile", cr.UserHandler.Profile, middlewareAut.IsLogin)
	router.POST("/api/profile", cr.UserHandler.UpdateProfile, middlewareAut.IsLogin)
	router.GET("/api/profile/recommend", cr.RecommendHandler.GetRecommendation, middlewareAut.SetTokenIfIsLogin)
	router.POST("/api/update/password", cr.UserHandler.UpdatePassword, middlewareAut.IsLogin)
	router.GET("/api/logout", cr.UserHandler.Logout, middlewareAut.IsLogin)
	router.GET("/api/homepage", cr.ProductHandler.GetAllProducts)
	router.GET("/api/product", cr.ProductHandler.GetProductById, middlewareAut.SetTokenIfIsLogin)
	router.POST("/api/product/add", cr.ProductHandler.AddProduct, middlewareAut.IsLogin)
	router.GET("/api/product/favorite/get", cr.ProductHandler.GetFavouriteProducts, middlewareAut.IsLogin)
	router.POST("/api/product/favorite/add", cr.ProductHandler.AddFavouriteProduct, middlewareAut.IsLogin)
	router.POST("/api/product/favorite/delete", cr.ProductHandler.DeleteFavouriteProduct, middlewareAut.IsLogin)
	router.POST("/api/sales/put", cr.ProductHandler.PutSalesForProduct, middlewareAut.IsLogin)
	router.GET("/api/sales", cr.ProductHandler.GetSalesProducts)
	router.GET("api/brands/get", cr.BrandHandler.GetBrands)
	router.GET("api/brand/products", cr.BrandHandler.GetProductsByBrand)
	router.POST("/api/upload/product", cr.ProductHandler.UploadProduct, middlewareAut.IsLogin)
	router.GET("/api/cart/get", cr.BasketHandler.GetBasket, middlewareAut.IsLogin)
	router.POST("/api/cart/put", cr.BasketHandler.PutInBasket, middlewareAut.IsLogin)
	router.POST("/api/cart/drop", cr.BasketHandler.DropBasket, middlewareAut.IsLogin)
	router.POST("/api/cart/delete", cr.BasketHandler.DeleteProduct, middlewareAut.IsLogin)
	router.POST("/api/cart/confirm", cr.OrderHandler.PutOrder, middlewareAut.IsLogin)
	router.POST("/api/cart/check", cr.OrderHandler.PutOrder, middlewareAut.IsLogin)
	router.POST("/api/category/create", cr.CategoryHandler.CreateCategory)
	router.GET("/api/category", cr.CategoryHandler.GetCategories)
	router.GET("/api/category/:name", cr.CategoryHandler.GetCategoryProducts)
	router.GET("/api/profile/orders", cr.OrderHandler.GetAllOrders, middlewareAut.IsLogin)
	router.POST("/api/review/add", cr.ReviewHandler.AddReview, middlewareAut.IsLogin)
	router.POST("/api/review/update", cr.ReviewHandler.UpdateReview, middlewareAut.IsLogin)
	router.DELETE("/api/review/delete", cr.ReviewHandler.DeleteReview, middlewareAut.IsLogin)
	router.GET("/api/reviews", cr.ReviewHandler.GetReviewsByProductId)
	router.GET("/api/user/reviews", cr.ReviewHandler.GetReviewsByUser)
	router.GET("/api/search/suggest", cr.SearchHandler.GetSuggests)
	router.GET("/api/search", cr.SearchHandler.GetSearchResults)
	router.GET("/api/metrics", echo.WrapHandler(promhttp.Handler()))
}
