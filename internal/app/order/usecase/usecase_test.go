package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_order "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/mocks"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestOrderUseCase_PutOrder(t *testing.T) {
	type mockBehaviorRepository func(s *mock_order.MockRepository, order models.Order)

	products := []models.OrderProducts {
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
		OrderId: 1,
		UserId: 3,
		Date:     "28-10-2021 03:03:59",
		Address:  address,
		Cost:     1000.99,
		Status:   "new",
		Products: products,
	}

	productPrices := []models.ProductPrice{
		{
			Id: 1,
			Price: 1000.99,
		},
	}

	tests := []struct {
		name                   string
		order              models.Order
		mockBehaviorRepository mockBehaviorRepository
		expectedId             int
		expectedError          error
	}{
		{
			name : "correct",
			order : order,
			mockBehaviorRepository: func(s *mock_order.MockRepository, order models.Order) {
				s.EXPECT().PutOrder(order).Return(3, nil)
				s.EXPECT().SelectPrices(order.Products).Return(productPrices, nil)
			},
			expectedId : 3,
			expectedError: nil,
		},
		{
			name : "error put order",
			order : order,
			mockBehaviorRepository: func(s *mock_order.MockRepository, order models.Order) {
				s.EXPECT().PutOrder(order).Return(0, errors.New("new error"))
				s.EXPECT().SelectPrices(order.Products).Return(productPrices, nil)
			},
			expectedId : 0,
			expectedError: errors.New("new error"),
		},
		{
			name : "error select prices",
			order : order,
			mockBehaviorRepository: func(s *mock_order.MockRepository, order models.Order) {
				s.EXPECT().SelectPrices(order.Products).Return(nil, errors.New("new error"))
			},
			expectedId : 0,
			expectedError: errors.New("new error"),
		},

	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			newOrderRepo := mock_order.NewMockRepository(c)

			tt.mockBehaviorRepository(newOrderRepo, tt.order)

			useCase := NewOrderUseCase(newOrderRepo)

			orderId, _, err := useCase.PutOrder(tt.order)

			assert.Equal(t, orderId, tt.expectedId)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}
