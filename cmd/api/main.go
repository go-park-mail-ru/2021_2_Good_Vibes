package main

//тут надо какой-то порядок с неймингами навести
import (
	"fmt"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configMiddleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configRouting"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configValidator"
	basketHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/delivery/http"
	basketUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/usecase"
	brandHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands/delivery/http"
	brandRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands/repository/posgresql"
	brandUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/brands/usecase"
	categoryHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/delivery/http"
	categoryRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/repository/posgresql"
	categoryUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/metrics"
	notification "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/notifications"
	notificationRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/notifications/repository/postgres"
	notificationUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/notifications/usecase"
	orderHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/delivery/http"
	orderUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/usecase"
	productHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	productRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/repository/postgresql"
	productUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/usecase"
	http3 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation/delivery/http"
	postgresql2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/recommendation/usecase"
	reviewHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/delivery/http"
	reviewRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/repository/postgresql"
	reviewUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/usecase"
	searchHandlerHttp "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search/delivery/http"
	searchRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search/repository/postgresql"
	searchUseCase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/search/usecase"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/manager"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	http2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/delivery/http"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/repository/postgresql"
	userUsecase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/usecase"
	orderRepoPostgres "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/repository/postgresql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var (
	router  = echo.New()
	storage user.Repository
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	sessionManager, err := manager.NewTokenManager(configApp.ConfigApp.SecretKey)
	if err != nil {
		logger.CustomLogger.LogrusLoggerHandler.Fatal(errors.BAD_INIT_SECRET_KEY)
	}

	os.Setenv("AWS_ACCESS_KEY", configApp.ConfigApp.AwsAccessKey)
	os.Setenv("AWS_SECRET_KEY", configApp.ConfigApp.AwsSecretKey)
	os.Setenv("DATABASE_URL", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		configApp.ConfigApp.DataBase.User, configApp.ConfigApp.DataBase.Password,
		configApp.ConfigApp.DataBase.Host, configApp.ConfigApp.DataBase.Port,
		configApp.ConfigApp.DataBase.DBName))

	//------------------user--------------------
	dbConn, dbErr := postgre.GetPostgres()

	storage, err = postgresql.NewStorageUserDB(dbConn, dbErr)
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	authGrpcConn, err := grpc.Dial(
		"localhost:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer authGrpcConn.Close()
	userHandler := http2.NewLoginHandler(userUsecase.NewUsecase(authGrpcConn, storage, impl.NewHasherBCrypt(bcrypt.DefaultCost)),
		sessionManager)

	//------------------product--------------------
	storageProd, err := productRepoPostgres.NewStorageProductsDB(dbConn, dbErr)
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	productHandler := productHandlerHttp.NewProductHandler(productUseCase.NewProductUsecase(storageProd),
		sessionManager)

	//------------------order--------------------
	orderGrpcConn, err := grpc.Dial(
		"localhost:8083",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	orderHandler := orderHandlerHttp.NewOrderHandler(orderUseCase.NewOrderUseCase(orderGrpcConn),
		sessionManager)

	//------------------basket--------------------
	BasketGrpcConn, err := grpc.Dial(
		"localhost:8082",
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	basketHandler := basketHandlerHttp.NewBasketHandler(basketUseCase.NewBasketUseCase(BasketGrpcConn),
		sessionManager)

	//------------------search--------------------
	storageSearch, err := searchRepoPostgres.NewSearchRepository(dbConn, dbErr)
	if err != nil {
		log.Fatal(err)
	}
	searchHandler := searchHandlerHttp.NewSearchHandler(searchUseCase.NewSearchUseCase(storageSearch))

	//------------------category--------------------
	storageCategory, err := categoryRepoPostgres.NewStorageCategoryDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}
	categoryHandler := categoryHandlerHttp.NewCategoryHandler(categoryUseCase.NewCategoryUseCase(storageCategory,
		storageProd))

	//-----------------brands----------------------
	storageBrand, err := brandRepoPostgres.NewStorageBrandDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}
	brandHandler := brandHandlerHttp.NewBrandHandler(brandUseCase.NewBrandUseCase(storageBrand, storageProd), sessionManager)

	//------------------reviews--------------------
	storageReview, err := reviewRepoPostgres.NewReviewRepository(dbConn, dbErr)
	if err != nil {
		panic(err)
	}
	reviewHandler := reviewHandlerHttp.NewReviewHandler(reviewUseCase.NewReviewUseCase(storageReview, storage),
		sessionManager)

	//------------------recommendation---------------
	storageReccomend, err := postgresql2.NewRecommendationRepository(dbConn, dbErr)
	if err != nil {
		panic(err)
	}
	recommendationHandler := http3.NewRecommendHandler(usecase.NewRecommendationUseCase(storageReccomend, storageProd), sessionManager)

	//------------------notification---------------
	notifyRepository, err := notificationRepoPostgres.NewStorageNotifyDB(dbConn, dbErr)
	if err != nil {
		panic(err)
	}

	orderRepository, err := orderRepoPostgres.NewOrderRepository(dbConn, dbErr)
	if err != nil {
		panic(err)
	}
	notifyUseCase := notificationUseCase.NewNotifyUseCase(notifyRepository, storage, orderRepository)
	notifyD := notification.NewNotifier(notifyUseCase)

	err = notifyD.Run()
	if err != nil {
		panic(err)
	}

	m, err := metrics.CreateNewMetric("main")
	if err != nil {
		panic(err)
	}

	router.Use(m.CollectMetrics)

	serverRouting := configRouting.ServerConfigRouting{
		ProductHandler:   productHandler,
		UserHandler:      userHandler,
		OrderHandler:     orderHandler,
		BasketHandler:    basketHandler,
		CategoryHandler:  categoryHandler,
		ReviewHandler:    reviewHandler,
		SearchHandler:    searchHandler,
		RecommendHandler: recommendationHandler,
		BrandHandler: brandHandler,
	}

	serverRouting.ConfigRouting(router)
	configValidator.ConfigValidator(router)

	configMiddleware.ConfigMiddleware(router)
	if err := router.Start(configApp.ConfigApp.MainConfig.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
