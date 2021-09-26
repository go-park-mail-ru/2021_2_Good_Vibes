package main

//тут надо какой-то порядок с неймингами навести
import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	storage_prod_handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/handler"
	storage_prod_useCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/storage"
	storage_prod_impl "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/storage/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/handler"
	middleware_user "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/middleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user/impl"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

var (
	router      = echo.New()
	storage     storage_user.UserUseCase
	storageProd storage_prod_useCase.UseCase
)


func main() {
	os.Setenv("DATABASE_URL", "postgres://lida:123@localhost:5432/mydb")

	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	storage, err = impl.NewStorageUserDB()
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	userHandler := handler.NewLoginHandler(&storage)
	router.Validator = &handler.CustomValidator{Validator: validator.New()}

	storageProd = storage_prod_impl.NewStorageProductsMemory()
	storageProd.AddProduct(product.Product{"1.jpg", "cat1", 10000})
	storageProd.AddProduct(product.Product{"2.jpg", "cat2", 10000})
	storageProd.AddProduct(product.Product{"3.jpg", "cat3", 10000})
	storageProd.AddProduct(product.Product{"4.jpg", "cat4", 10000})
	storageProd.AddProduct(product.Product{"5.jpg", "cat5", 10000})
	storageProd.AddProduct(product.Product{"6.jpg", "dog", 100001})

	productHandler := storage_prod_handler.NewProductHandler(storageProd)
	//для этого инит какой то надо придумать
	router.Static("/", "static")
	router.POST("/login", userHandler.Login)
	router.POST("/signup", userHandler.SignUp)
	router.GET("/profile", profile, middleware_user.IsLogin)
	router.GET("/homepage", productHandler.GetAllProducts)
	router.GET("/logout", userHandler.Logout, middleware_user.IsLogin)

	if err := router.Start(configApp.ConfigApp.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
