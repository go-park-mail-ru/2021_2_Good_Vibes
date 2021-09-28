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


	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	os.Setenv("DATABASE_URL", configApp.ConfigApp.DataBaseURL)

	storage, err = impl.NewStorageUserMemory()
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	userHandler := handler.NewLoginHandler(&storage)
	val := validator.New()
	val.RegisterValidation("customPassword", handler.Password)
	router.Validator = &handler.CustomValidator{Validator: val}

	storageProd = storage_prod_impl.NewStorageProductsMemory()
	storageProd.AddProduct(product.Product{1,"images/shoe2.png", "Кроссовки adidas голубые", 250, 4})
	storageProd.AddProduct(product.Product{2,"images/phone2.png", "Смартфон", 10000, 2.5})
	storageProd.AddProduct(product.Product{3,"images/shirt1.png", "Кофта мужская", 10000, 2.5})
	storageProd.AddProduct(product.Product{4,"images/smartphone.png", "Смартфон чёрный цвет", 10000, 2.5})
	storageProd.AddProduct(product.Product{5,"images/shirt4.png", "Кофта мужская", 10000, 2.5})
	storageProd.AddProduct(product.Product{6,"images/shoe5.png", "Кеды adidas желтые", 10000, 2.5})
	storageProd.AddProduct(product.Product{7,"images/phone3.png", "Смартфон поддержанный", 10000, 2.5})
	storageProd.AddProduct(product.Product{8,"images/shoe1.png", "Кроссовки adidas красные", 10000, 2.5})
	storageProd.AddProduct(product.Product{9,"images/shoe3.png", "Кроссовки adidas черные", 10000, 2.5})

	productHandler := storage_prod_handler.NewProductHandler(storageProd)
	//для этого инит какой то надо придумать
	router.Static("/", "static")
	router.POST("/login", userHandler.Login)
	router.POST("/signup", userHandler.SignUp)
	router.GET("/profile", profile, middleware_user.IsLogin)
	router.GET("/homepage", productHandler.GetAllProducts)
	router.GET("/product", productHandler.GetProductById)
	router.GET("/logout", userHandler.Logout, middleware_user.IsLogin)

	if err := router.Start(configApp.ConfigApp.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
