package models

import proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"

type OrderProducts struct {
	OrderId        int     `json:"order_id,omitempty"`
	ProductId      int     `json:"product_id" validate:"required"`
	Number         int     `json:"number" validate:"required"`
	Price          float64 `json:"price,omitempty"`
	PriceWithPromo float64 `json:"price_with_promo,omitempty"`
}

type Address struct {
	Country string `json:"country" validate:"required"`
	Region  string `json:"region" validate:"required"`
	City    string `json:"city" validate:"required"`
	Street  string `json:"street" validate:"required"`
	House   string `json:"house" validate:"required"`
	Flat    string `json:"flat" validate:"required"`
	Index   string `json:"index" validate:"required"`
}

type Order struct {
	OrderId   int             `json:"order_id,omitempty"`
	UserId    int             `json:"user_id,omitempty"`
	Date      string          `json:"date,omitempty"`
	Address   Address         `json:"address,omitempty"`
	Cost      float64         `json:"cost,omitempty"`
	CostWithPromo float64	  `json:"cost_with_promo, omitempty"`
	Status    string          `json:"status,omitempty"`
	Products  []OrderProducts `json:"products" validate:"required"`
	Promocode string          `json:"promocode,omitempty"`
	Email     string          `json:"email"`
}

func GrpcAddressToModel(grpcData *proto.Address) Address {
	return Address{
		Country: grpcData.GetCountry(),
		Region:  grpcData.GetRegion(),
		City:    grpcData.GetCity(),
		Street:  grpcData.GetStreet(),
		House:   grpcData.GetHouse(),
		Flat:    grpcData.GetFlat(),
		Index:   grpcData.GetIndex(),
	}
}

func ModelAddressToGrpc(model Address) *proto.Address {
	return &proto.Address{
		Country: model.Country,
		Region:  model.Region,
		City:    model.City,
		Street:  model.Street,
		House:   model.House,
		Flat:    model.Flat,
		Index:   model.Index,
	}
}

func GrpcOrderProductsToModel(grpcData *proto.OrderProducts) OrderProducts {
	return OrderProducts{
		OrderId:        int(grpcData.GetOrderId()),
		ProductId:      int(grpcData.GetProductId()),
		Number:         int(grpcData.GetNumber()),
		Price:          float64(grpcData.GetPrice()),
		PriceWithPromo: float64(grpcData.GetPriceWithPromo()),
	}
}

func ModelOrderProductsToGrpc(model OrderProducts) *proto.OrderProducts {
	return &proto.OrderProducts{
		OrderId:        int64(model.OrderId),
		ProductId:      int64(model.ProductId),
		Number:         int64(model.Number),
		Price:          float32(model.Price),
		PriceWithPromo: float32(model.PriceWithPromo),
	}
}

func GrpcOrderToModel(grpcData *proto.Order) Order {
	var productsModel []OrderProducts
	for _, element := range grpcData.GetProducts() {
		productsModel = append(productsModel, GrpcOrderProductsToModel(element))
	}

	return Order{
		OrderId:   int(grpcData.GetOrderId()),
		UserId:    int(grpcData.GetUserId()),
		Date:      grpcData.GetDate(),
		Address:   GrpcAddressToModel(grpcData.GetAddress()),
		Cost:      float64(grpcData.GetCost()),
		CostWithPromo: float64(grpcData.GetCostWithPromo()),
		Products:  productsModel,
		Promocode: grpcData.GetPromocode(),
		Email:     grpcData.GetEmail(),
		Status:    grpcData.GetStatus(),
	}
}

func ModelOrderToGrpc(model Order) *proto.Order {
	productsProto := []*proto.OrderProducts{}
	for _, element := range model.Products {
		productsProto = append(productsProto, ModelOrderProductsToGrpc(element))
	}

	return &proto.Order{
		OrderId:   int64(model.OrderId),
		UserId:    int64(model.UserId),
		Date:      model.Date,
		Address:   ModelAddressToGrpc(model.Address),
		Cost:      float32(model.Cost),
		CostWithPromo: float32(model.CostWithPromo),
		Products:  productsProto,
		Promocode: model.Promocode,
		Status:    model.Status,
		Email:     model.Email,
	}
}
