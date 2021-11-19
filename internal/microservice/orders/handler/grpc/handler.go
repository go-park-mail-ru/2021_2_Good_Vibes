package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/orders"
)

type grpcOrderHandler struct {
	orderUseCase orders.UseCase
	proto.UnimplementedOrderServiceServer
}

func NewGrpcOrderHandler(userUseCase orders.UseCase) *grpcOrderHandler {
	return &grpcOrderHandler{
		orderUseCase: userUseCase,
	}
}

func (handler *grpcOrderHandler) PutOrder(ctx context.Context, order *proto.Order) (*proto.OrderCost, error) {
	orderId, cost, err := handler.orderUseCase.PutOrder(models.GrpcOrderToModel(order))
	if err != nil {
		return nil, err
	}
	return &proto.OrderCost{OrderId: int64(orderId), Cost: float32(cost)}, nil
}

func (handler *grpcOrderHandler) GetAllOrders (ctx context.Context, userId *proto.UserIdOrder) (*proto.ArrayOrders, error) {
	ordersModel, err := handler.orderUseCase.GetAllOrders(int(userId.GetUserId()))
	if err != nil {
		return nil, err
	}

	ordersGrpc := []*proto.Order{}
	for _, element := range ordersModel {
		ordersGrpc = append(ordersGrpc, models.ModelOrderToGrpc(element))
	}
	return &proto.ArrayOrders{
		Orders: ordersGrpc,
	}, nil
}