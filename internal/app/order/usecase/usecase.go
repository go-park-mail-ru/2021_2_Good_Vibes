package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"
	"google.golang.org/grpc"
	"time"
)

type UseCase struct {
	orderServiceClient proto.OrderServiceClient
}

func NewOrderUseCase(conn *grpc.ClientConn) *UseCase {
	c := proto.NewOrderServiceClient(conn)

	return &UseCase{
		orderServiceClient: c,
	}
}

func (uc *UseCase) PutOrder(order models.Order) (int, float64, error) {
	order.Date = time.Now().Format(time.RFC3339)
	orderCost, err := uc.orderServiceClient.PutOrder(context.Background(), models.ModelOrderToGrpc(order))
	if err != nil {
		return 0, 0, err
	}
	return int(orderCost.GetOrderId()), float64(orderCost.GetCost()), nil
}

func (uc *UseCase) GetAllOrders(user int) ([]models.Order, error) {
	ordersGrpc, err := uc.orderServiceClient.GetAllOrders(context.Background(), &proto.UserIdOrder{UserId: int64(user)})
	if err != nil {
		return nil, err
	}

	var ordersModel []models.Order
	for _, element := range ordersGrpc.GetOrders() {
		ordersModel = append(ordersModel, models.GrpcOrderToModel(element))
	}
	return ordersModel, nil
}
