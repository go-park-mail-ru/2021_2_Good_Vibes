package main

import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/auth"
	handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/handler/grpc"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/usecase"
	_ "github.com/jackc/pgx/stdlib"
	"golang.org/x/crypto/bcrypt"
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

	storage, err := postgresql.NewStorageUserDB(postgre.GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	hasher := impl.NewHasherBCrypt(bcrypt.DefaultCost)
	userUseCase := usecase.NewUsecase(storage, hasher)
	handler := handler.NewGrpcUserHandler(userUseCase)

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("can't listen auth microservice")
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
