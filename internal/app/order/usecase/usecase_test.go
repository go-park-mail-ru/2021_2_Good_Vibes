package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_order "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/mocks"
	proto "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/order"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
)

var errGetting = errors.New("grpc error")

func TestUseCase_GetAllOrders(t *testing.T) {
	type mockBehaviorClientGRPC func(s *mock_order.MockOrderServiceClient,
		ctx context.Context, product *proto.UserIdOrder)

	products := []models.OrderProducts{
		{
			OrderId:   1,
			ProductId: 10,
			Number:    2,
		},
		{
			OrderId:   1,
			ProductId: 1,
			Number:    1,
		},
		{
			OrderId:   1,
			ProductId: 3,
			Number:    4,
		},
	}

	address := models.Address{
		Country: "Russia",
		Region:  "Moscow",
		City:    "Moscow",
		Street:  "Izmailovskiy prospect",
		House:   "73B",
		Flat:    "44",
		Index:   "109834",
	}

	order := models.Order{
		OrderId:  1,
		UserId:   3,
		Date:     "2021-11-23T00:33:46+03:00",
		Address:  address,
		Cost:     50000.00,
		Status:   "",
		Products: products,
	}
	var orders = []models.Order{
		order,
	}

	ordersGrpc := []*proto.Order{}
	for _, element := range orders {
		ordersGrpc = append(ordersGrpc, models.ModelOrderToGrpc(element))
	}

	testTable := []struct {
		name                   string
		inputData              int
		mockBehaviorClientGRPC mockBehaviorClientGRPC
		expectedData           []models.Order
		expectedError          error
	}{
		{
			name:      "ok",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_order.MockOrderServiceClient, ctx context.Context, product *proto.UserIdOrder) {
				s.EXPECT().GetAllOrders(ctx, product).Return(&proto.ArrayOrders{
					Orders: ordersGrpc,
				}, nil)
			},
			expectedData:  orders,
			expectedError: nil,
		},
		{
			name:      "fail",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_order.MockOrderServiceClient, ctx context.Context, product *proto.UserIdOrder) {
				s.EXPECT().GetAllOrders(ctx, product).Return(&proto.ArrayOrders{
					Orders: ordersGrpc,
				}, errGetting)
			},
			expectedData:  nil,
			expectedError: errGetting,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockClientGRPC := mock_order.NewMockOrderServiceClient(c)
			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
				&proto.UserIdOrder{UserId: int64(testCase.inputData)})

			useCase := UseCase{orderServiceClient: mockClientGRPC}

			data, err := useCase.GetAllOrders(testCase.inputData)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedData, data)
		})
	}
}

// тут со временем проблемы
//func TestUseCase_PutOrder(t *testing.T) {
//	type mockBehaviorClientGRPC func (s *mock_order.MockOrderServiceClient,
//		ctx context.Context, order *proto.Order)
//
//	products := []models.OrderProducts{
//		{
//			OrderId:   1,
//			ProductId: 10,
//			Number:    2,
//		},
//		{
//			OrderId:   1,
//			ProductId: 1,
//			Number:    1,
//		},
//		{
//			OrderId:   1,
//			ProductId: 3,
//			Number:    4,
//		},
//	}
//
//	address := models.Address{
//		Country: "Russia",
//		Region:  "Moscow",
//		City:    "Moscow",
//		Street:  "Izmailovskiy prospect",
//		House:   "73B",
//		Flat:    "44",
//		Index:   "109834",
//	}
//
//	order := models.Order{
//		OrderId:  1,
//		UserId:   3,
//		Date:     "2021-11-23T00:33:46+03:00",
//		Address:  address,
//		Cost:     50000.00,
//		Status:   "",
//		Products: products,
//	}
//
//	testTable := []struct{
//		name string
//		mockBehaviorClientGRPC mockBehaviorClientGRPC
//		expectedData1 int
//		expectedData2 float64
//		expectedError error
//	} {
//		{
//			name: "ok",
//			mockBehaviorClientGRPC: func(s *mock_order.MockOrderServiceClient, ctx context.Context, order *proto.Order) {
//				s.EXPECT().PutOrder(ctx, order).Return(&proto.OrderCost{
//					OrderId: order.OrderId,
//					Cost: order.Cost,
//				}, nil)
//			},
//			expectedData1: order.OrderId,
//			expectedData2: order.Cost,
//			expectedError: nil,
//		},
//		{
//			name: "fail",
//			mockBehaviorClientGRPC: func(s *mock_order.MockOrderServiceClient, ctx context.Context, order *proto.Order) {
//				s.EXPECT().PutOrder(ctx, order).Return(&proto.OrderCost{
//					OrderId: order.OrderId,
//					Cost: order.Cost,
//				}, errGetting)
//			},
//			expectedData1: 0,
//			expectedData2: 0,
//			expectedError: errGetting,
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockClientGRPC := mock_order.NewMockOrderServiceClient(c)
//			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
//				models.ModelOrderToGrpc(order))
//
//			useCase := UseCase{orderServiceClient: mockClientGRPC}
//
//			data1, data2, err := useCase.PutOrder(order)
//
//			assert.Equal(t, testCase.expectedError, err)
//			assert.Equal(t, testCase.expectedData1, data1)
//			assert.Equal(t, testCase.expectedData2, data2)
//		})
//	}
//}
