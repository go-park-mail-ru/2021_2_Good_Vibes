package main

//тут надо какой-то порядок с неймингами навести
import (
	"database/sql"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configRouting"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/config/configValidator"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	storage_prod_handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
	storage_prod_impl "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/repository/memory"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/usecase"
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
	router      = echo.New()
	storage     user.Repository
	storageProd product.Repository
)

func main() {

	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	os.Setenv("DATABASE_URL", configApp.ConfigApp.DataBaseURL)

	//storage, err = impl.NewStorageUserDB(GetPostgres())
	storage, err = postgresql.NewStorageUserDB(GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	userUс := userUsecase.NewUsecase(storage)

	//storageProd, err = storage_prod_impl.NewStorageProductsDB(GetPostgres())
	storageProd, err = storage_prod_impl.NewStorageProductsMemory()
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}
	productUc := usecase.NewProductUsecase(storageProd)

	productUc.AddProduct(models.Product{Id: 1, Image: "images/shoe2.png", Name: "Кроссовки adidas голубые", Price: 250, Rating: 4, Category: "SNICKERS_ADIDAS_MEN"})
	productUc.AddProduct(models.Product{2, "images/phone2.png", "Смартфон", 10000, 2.5, "PHONES"})
	productUc.AddProduct(models.Product{3, "images/shirt1.png", "Кофта мужская", 10000, 2.5, "CLOTHES_UP_MEN"})
	productUc.AddProduct(models.Product{4, "images/smartphone.png", "Смартфон чёрный цвет", 10000, 2.5, "PHONES"})
	productUc.AddProduct(models.Product{5, "images/shirt4.png", "Кофта мужская", 10000, 2.5, "CLOTHES_UP_MEN"})
	productUc.AddProduct(models.Product{6, "images/shoe5.png", "Кеды adidas желтые", 10000, 2.5, "SNICKERS_ADIDAS_MEN"})
	productUc.AddProduct(models.Product{7, "images/phone3.png", "Смартфон поддержанный", 10000, 2.5, "PHONES"})
	productUc.AddProduct(models.Product{8, "images/shoe1.png", "Кроссовки adidas красные", 10000, 2.5, "SNICKERS_ADIDAS_MEN"})
	productUc.AddProduct(models.Product{9, "images/shoe3.png", "Кроссовки adidas черные", 10000, 2.5, "SNICKERS_ADIDAS_MEN"})

	productHandler := storage_prod_handler.NewProductHandler(productUc)

	userHandler := http2.NewLoginHandler(userUс)

	serverRouting := configRouting.ServerConfigRouting{ProductHandler: productHandler,
		UserHandler: userHandler}
	serverRouting.ConfigRouting(router)
	configValidator.ConfigValidator(router)

	if err := router.Start(configApp.ConfigApp.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func GetPostgres() (*sql.DB, error) {
	dsn := "user=lida dbname=mydb password=123 host=127.0.0.1 port=5432 sslmode=disable"
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
