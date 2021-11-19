package main

//тут надо какой-то порядок с неймингами навести
import (
	"database/sql"
	"fmt"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configMiddleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configRouting"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configValidator"
	basketHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	basketUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	orderHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	orderUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	productHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	productRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/manager"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

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
	storageCategory category.Repository
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig("../..")
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

	authGrpcConn, err := grpc.Dial(
		"localhost:8081",
				grpc.WithInsecure(),
		)
	if err != nil {
		log.Fatal(err)
	}
	defer authGrpcConn.Close()
	userUс := userUsecase.NewUsecase(authGrpcConn, storage, hasher)


	storageProd, err = productRepoPostgres.NewStorageProductsDB(GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	productUc := productUseCase.NewProductUsecase(storageProd)


	orderGrpcConn, err := grpc.Dial(
		"localhost:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	orderUc := orderUseCase.NewOrderUseCase(orderGrpcConn)

	sessionManager, err := manager.NewTokenManager(configApp.ConfigApp.SecretKey)
	if err != nil {
		logger.CustomLogger.LogrusLoggerHandler.Fatal(errors.BAD_INIT_SECRET_KEY)
	}
	userHandler := http2.NewLoginHandler(userUс, sessionManager)

	BasketGrpcConn, err := grpc.Dial(
		"localhost:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	basketUc := basketUseCase.NewBasketUseCase(BasketGrpcConn)
	basketHandler := basketHandlerHttp.NewBasketHandler(basketUc, sessionManager)

	storageCategory, err := categoryRepoPostgres.NewStorageCategoryDB(GetPostgres())
	if err != nil {
		panic(err)
	}

	categoryUc := categoryUseCase.NewCategoryUseCase(storageCategory, storageProd)

	productHandler := productHandlerHttp.NewProductHandler(productUc, sessionManager)
	orderHandler := orderHandlerHttp.NewOrderHandler(orderUc, sessionManager)
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
