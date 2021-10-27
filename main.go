package main

//тут надо какой-то порядок с неймингами навести
import (
	"database/sql"
	"fmt"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configMiddleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configRouting"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configValidator"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket"
	basketHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	basketRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/repository/postgresql"
	basketUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order"
	orderHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	orderRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/repository/postgresql"
	orderUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	productHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	productRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/manager"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category"
	categoryHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/delivery/http"
	categoryRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/repository/posgresql"
	categoryUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/usecase"

	productUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	http2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/usecase"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

var (
	router          = echo.New()
	storage         user.Repository
	storageProd     product.Repository
	storageOrder    order.Repository
	storageBasket   basket.Repository
	storageCategory category.Repository
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	os.Setenv("AWS_ACCESS_KEY", configApp.ConfigApp.AwsAccessKey)
	os.Setenv("AWS_SECRET_KEY", configApp.ConfigApp.AwsSecretKey)
	os.Setenv("DATABASE_URL", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configApp.ConfigApp.DataBase.User, configApp.ConfigApp.DataBase.Password,
		configApp.ConfigApp.DataBase.Host, configApp.ConfigApp.DataBase.Port,
		configApp.ConfigApp.DataBase.DBName))

	storage, err = postgresql.NewStorageUserDB(GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	hasher := impl.NewHasherBCrypt(bcrypt.DefaultCost)
	userUс := userUsecase.NewUsecase(storage, hasher)

	storageProd, err = productRepoPostgres.NewStorageProductsDB(GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	productUc := productUseCase.NewProductUsecase(storageProd)

	productUc.AddProduct(models.Product{Image: "images/shoe2.png", Name: "Кроссовки adidas голубые", Price: 250, Rating: 4, Category: "SNICKERS_ADIDAS_MEN", CountInStock: 100, Description: "Крутые adidas кроссовки"})
	productUc.AddProduct(models.Product{Image: "images/phone2.png", Name: "Смартфон", Price: 10000, Rating: 2.5, Category: "PHONES", CountInStock: 100, Description: "Крутой новый смартфон"})
	productUc.AddProduct(models.Product{Image: "images/shirt1.png", Name: "Кофта мужская", Price: 10000, Rating: 2.5, Category: "CLOTHES_UP_MEN", CountInStock: 100, Description: "Крутая мужская кофта"})
	productUc.AddProduct(models.Product{Image: "images/smartphone.png", Name: "Смартфон чёрный цвет", Price: 10000, Rating: 2.5, Category: "PHONES", CountInStock: 100, Description: "Крутой черный смартфон"})
	productUc.AddProduct(models.Product{Image: "images/shirt4.png", Name: "Кофта мужская", Price: 10000, Rating: 2.5, Category: "CLOTHES_UP_MEN", CountInStock: 100, Description: "Супер крутая красный смартфон"})
	productUc.AddProduct(models.Product{Image: "images/shoe5.png", Name: "Кеды adidas желтые", Price: 10000, Rating: 2.5, Category: "SNICKERS_ADIDAS_MEN", CountInStock: 100, Description: "Крутые желтые кроссовки"})
	productUc.AddProduct(models.Product{Image: "images/phone3.png", Name: "Смартфон поддержанный", Price: 10000, Rating: 2.5, Category: "PHONES", CountInStock: 100, Description: "Крутой поддержанный смартфон"})
	productUc.AddProduct(models.Product{Image: "images/shoe1.png", Name: "Кроссовки adidas красные", Price: 10000, Rating: 2.5, Category: "SNICKERS_ADIDAS_MEN", CountInStock: 100, Description: "Крутые красные кроссовки"})
	productUc.AddProduct(models.Product{Image: "images/shoe3.png", Name: "Кроссовки adidas черные", Price: 10000, Rating: 2.5, Category: "SNICKERS_ADIDAS_MEN", CountInStock: 100, Description: "Крутые черные кроссовкин"})

	storageOrder, err := orderRepoPostgres.NewOrderRepository(GetPostgres())
	if err != nil {
		panic(err)
	}

	orderUc := orderUseCase.NewOrderUseCase(storageOrder)

	storageBasket, err := basketRepoPostgres.NewBasketRepository(GetPostgres())
	basketUc := basketUseCase.NewBasketUseCase(storageBasket)
	basketHandler := basketHandlerHttp.NewBasketHandler(basketUc)

	storageCategory, err := categoryRepoPostgres.NewStorageCategoryDB(GetPostgres())
	if err != nil {
		panic(err)
	}

	categoryUc := categoryUseCase.NewCategoryUseCase(storageCategory, storageProd)

	sessionManager, err := manager.NewTokenManager(configApp.ConfigApp.SecretKey)
	{
		if err != nil {
			logger.CustomLogger.LogrusLoggerHandler.Fatal(errors.BAD_INIT_SECRET_KEY)
		}
	}

	productHandler := productHandlerHttp.NewProductHandler(productUc)
	userHandler := http2.NewLoginHandler(userUс, sessionManager)
	orderHandler := orderHandlerHttp.NewOrderHandler(orderUc)
	categoryHandler := categoryHandlerHttp.NewCategoryHandler(categoryUc)

	serverRouting := configRouting.ServerConfigRouting{
		ProductHandler:  productHandler,
		UserHandler:     userHandler,
		OrderHandler:    orderHandler,
		BasketHandler:   basketHandler,
		CategoryHandler: categoryHandler,
	}
	serverRouting.ConfigRouting(router)
	configValidator.ConfigValidator(router)

	configMiddleware.ConfigMiddleware(router)
	if err := router.Start(configApp.ConfigApp.MainConfig.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func GetPostgres() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		configApp.ConfigApp.DataBase.User, configApp.ConfigApp.DataBase.DBName,
		configApp.ConfigApp.DataBase.Password, configApp.ConfigApp.DataBase.Host,
		configApp.ConfigApp.DataBase.Port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}
