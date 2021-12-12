package main

import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/impl"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/auth"
	handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/handler/grpc"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth/usecase"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"
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

	storage, err := postgresql.NewStorageUserDB(postgre.GetPostgres())
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	handler := handler.NewGrpcUserHandler(usecase.NewUsecase(storage,
		impl.NewHasherBCrypt(bcrypt.DefaultCost)))

	m, err := metrics.CreateNewMetric("auth")
	if err != nil {
		log.Fatal("create metric error", err)
	}
	interceptor := metrics.NewInterceptor(m)

	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("can't listen auth microservice")
	}
	defer lis.Close()

	go func() {
		r := echo.New()
		r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		lis, err := net.Listen("tcp", "localhost:7000")
		if err != nil {
			log.Fatal(err)
			return
		}

		r.Listener = lis
		r.Start("")
	}()

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(), interceptor.Collect),
	)
	proto.RegisterAuthServiceServer(server, handler)
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
