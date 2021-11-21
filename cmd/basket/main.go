package main

import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/basket"
	handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/basket/handler/grpc"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/basket/repository/postgresql"
	usecase "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/basket/usecase"
	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig("/home/ubuntu/Ozon/2021_2_Good_Vibes")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	storage, err := postgresql.NewBasketRepository(postgre.GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	handler := handler.NewGrpcBasketHandler(usecase.NewBasketUseCase(storage))

	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatal("can't listen auth microservice")
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterBasketServiceServer(server, handler)
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
