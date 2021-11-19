package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/basket"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/basket"
	"google.golang.org/protobuf/types/known/emptypb"
)

type grpcBasketHandler struct {
	basketUseCase basket.UseCase
	proto.UnimplementedBasketServiceServer
}

func NewGrpcBasketHandler(userUseCase basket.UseCase) *grpcBasketHandler {
	return &grpcBasketHandler{
		basketUseCase: userUseCase,
	}
}

func (handler *grpcBasketHandler) PutInBasket(ctx context.Context, product *proto.BasketProduct) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, handler.basketUseCase.PutInBasket(models.GrpcBasketProductToModel(product))
}

func (handler *grpcBasketHandler) GetBasket(ctx context.Context, userId *proto.UserID) (*proto.ProductArray, error) {
	products, err := handler.basketUseCase.GetBasket(int(userId.GetUserId()))
	if err != nil {
		return nil, err
	}
	productsProtoArray := &proto.ProductArray{}
	for _, element := range products {
		productsProtoArray.Products = append(productsProtoArray.Products, models.ModelBasketProductToGrpc(element))
	}
	return productsProtoArray, nil
}

func (handler *grpcBasketHandler) DropBasket(ctx context.Context, userId *proto.UserID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, handler.basketUseCase.DropBasket(int(userId.GetUserId()))
}

func (handler *grpcBasketHandler) DeleteProduct(ctx context.Context, product *proto.BasketProduct) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, handler.basketUseCase.DeleteProduct(models.GrpcBasketProductToModel(product))
}
