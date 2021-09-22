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
)

var (
	router  = echo.New()
	storage storage_user.UserUseCase
)

func main() {
	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	storage = impl.NewStorageUserMemory()
	storage.AddUser(storage_user.User{"Misha", "qwerty@gmail.com", "1234"})
	storage.AddUser(storage_user.User{"Glasha", "qwerty@gmail.com", "1234"})
	storage.AddUser(storage_user.User{"Vova", "qwerty@gmail.com", "1234"})

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
