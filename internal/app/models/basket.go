package models

import proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/basket"

type BasketProduct struct {
	UserId    int `json:"user_id,omitempty"`
	ProductId int `json:"product_id"`
	Number    int `json:"number,omitempty"`
}

type BasketStorage struct {
	UserId int `json:"user_id"`
}

func ModelBasketProductToGrpc (model BasketProduct) *proto.BasketProduct{
	return &proto.BasketProduct{
		UserId: int64(model.UserId),
		ProductId: int64(model.ProductId),
		Number: int64(model.Number),
	}
}

func GrpcBasketProductToModel(grpcData *proto.BasketProduct) BasketProduct {
	return BasketProduct{
		UserId: int(grpcData.GetUserId()),
		ProductId: int(grpcData.GetProductId()),
		Number: int(grpcData.GetNumber()),
	}
}

func ModelBasketStorageToGrpc (model BasketStorage) *proto.UserID{
	return &proto.UserID{
		UserId: int64(model.UserId),
	}
}
