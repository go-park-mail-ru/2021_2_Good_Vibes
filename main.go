package main

import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/handler"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/middleware"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user/impl"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

var (
	router  = echo.New()
	storage storage_user.UserUseCase
)

func main() {
	os.Setenv("DATABASE_URL", "postgres://lida:sergeykust000@localhost:5432/mydb")

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

	router.Static("/", "static")
	router.POST("/login", userHandler.Login)
	router.POST("/signup", userHandler.SignUp)
	router.GET("/profile", profile, middleware.IsLogin)

	if err := router.Start(configApp.ConfigApp.ServerAddress); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

//пока просто для проверки middleware
func profile(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello world")
}
