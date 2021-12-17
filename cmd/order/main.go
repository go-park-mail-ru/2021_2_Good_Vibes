package main

import (
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/metrics"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/logger"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"
	userRepo "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/repository/postgresql"
	handler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/handler/grpc"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/repository/postgresql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders/usecase"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger.InitLogger()
	err := configApp.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	dbConn, dbErr := postgre.GetPostgres()
	storage, err := postgresql.NewOrderRepository(dbConn, dbErr)
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	userStorage, err := userRepo.NewStorageUserDB(dbConn, dbErr)
	if err != nil {
		log.Fatal("cannot connect data base", err)
	}

	handler := handler.NewGrpcOrderHandler(usecase.NewOrderUseCase(storage, userStorage))
	m, err := metrics.CreateNewMetric("order")
	if err != nil {
		log.Fatal("create metric error", err)
	}
	interceptor := metrics.NewInterceptor(m)

	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatal("can't listen auth microservice")
	}
	defer lis.Close()

	go func() {
		r := echo.New()
		r.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		lis, err := net.Listen("tcp", "localhost:7002")
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
	proto.RegisterOrderServiceServer(server, handler)
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
