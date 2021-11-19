package main

import (
	"database/sql"
	"fmt"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"
	handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/handler/grpc"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/repository/postgresql"
	usecase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/usecase"
	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	storage, err := postgresql.NewOrderRepository(GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	orderUseCase := usecase.NewOrderUseCase(storage)
	handler := handler.NewGrpcOrderHandler(orderUseCase)

	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatal("can't listen auth microservice")
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterOrderServiceServer(server, handler)

	if err := server.Serve(lis); err != nil {
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
