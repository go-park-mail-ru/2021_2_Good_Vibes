package usecase

import (
	"context"
	"errors"
	mock_basket "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/mocks"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/proto/basket"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
)

var errGetting = errors.New("grpc error")

func TestUseCase_PutInBasket(t *testing.T) {
	type mockBehaviorClientGRPC func (s *mock_basket.MockBasketServiceClient,
		ctx context.Context, product *basket.BasketProduct)

	testTable := []struct{
		name string
		inputData models.BasketProduct
		mockBehaviorClientGRPC mockBehaviorClientGRPC
		expectedError error
	}{
		{
			name: "OK",
			inputData: models.BasketProduct{
				UserId: 1,
				ProductId: 1,
				Number: 1,
			},
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, product *basket.BasketProduct) {
				s.EXPECT().PutInBasket(ctx, product).Return(nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "FAIL",
			inputData: models.BasketProduct{
				UserId: 1,
				ProductId: 1,
				Number: 1,
			},
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, product *basket.BasketProduct) {
				s.EXPECT().PutInBasket(ctx, product).Return(nil, errGetting)
			},
			expectedError: errGetting,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockClientGRPC := mock_basket.NewMockBasketServiceClient(c)
			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
				models.ModelBasketProductToGrpc(testCase.inputData))

			useCase := UseCase{
				basketServiceClient: mockClientGRPC,
			}

			assert.Equal(t, testCase.expectedError, useCase.PutInBasket(testCase.inputData))
		})
	}
}

func TestUseCase_GetBasket(t *testing.T) {
	type mockBehaviorClientGRPC func (s *mock_basket.MockBasketServiceClient,
		ctx context.Context,  id *basket.UserID)

	gettingData := []models.BasketProduct{
		{
			ProductId: 1,
			UserId:    1,
			Number:    1,
		},
		{
			ProductId: 2,
			UserId:    2,
			Number:    2,
		},
	}

	gettingDataGRPC := &basket.ProductArray{}
	for _, element := range gettingData {
		gettingDataGRPC.Products = append(gettingDataGRPC.Products,  models.ModelBasketProductToGrpc(element))
	}
	testTable := []struct{
		name string
		inputData int
		mockBehaviorClientGRPC mockBehaviorClientGRPC
		expectedData []models.BasketProduct
		expectedError error
	}{
		{
			name: "OK",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, id *basket.UserID) {
				s.EXPECT().GetBasket(ctx, id).Return(gettingDataGRPC, nil)
			},
			expectedData: gettingData,
			expectedError: nil,
		},
		{
			name: "FAIL",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, id *basket.UserID) {
				s.EXPECT().GetBasket(ctx, id).Return(nil, errGetting)
			},
			expectedData: nil,
			expectedError: errGetting,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockClientGRPC := mock_basket.NewMockBasketServiceClient(c)
			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
				models.ModelBasketStorageToGrpc(models.BasketStorage{UserId:
					testCase.inputData}))

			useCase := UseCase{
				basketServiceClient: mockClientGRPC,
			}
			gettingData, err := useCase.GetBasket(testCase.inputData)
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedData, gettingData)
		})
	}
}

func TestUseCase_DropBasket(t *testing.T) {
	type mockBehaviorClientGRPC func (s *mock_basket.MockBasketServiceClient,
		ctx context.Context,  id *basket.UserID)
	testTable := []struct{
		name string
		inputData int
		mockBehaviorClientGRPC mockBehaviorClientGRPC
		expectedError error
	}{
		{
			name: "OK",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, id *basket.UserID) {
				s.EXPECT().DropBasket(ctx, id).Return( nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "FAIL",
			inputData: 1,
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, id *basket.UserID) {
				s.EXPECT().DropBasket(ctx, id).Return(nil, errGetting)
			},
			expectedError: errGetting,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockClientGRPC := mock_basket.NewMockBasketServiceClient(c)
			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
				models.ModelBasketStorageToGrpc(models.BasketStorage{UserId:
				testCase.inputData}))

			useCase := UseCase{
				basketServiceClient: mockClientGRPC,
			}

			assert.Equal(t, testCase.expectedError, useCase.DropBasket(testCase.inputData))
		})
	}
}

func TestUseCase_DeleteProduct(t *testing.T) {
	type mockBehaviorClientGRPC func (s *mock_basket.MockBasketServiceClient,
		ctx context.Context,  product *basket.BasketProduct)

	testTable := []struct{
		name string
		inputData models.BasketProduct
		mockBehaviorClientGRPC mockBehaviorClientGRPC
		expectedError error
	}{
		{
			name: "OK",
			inputData: models.BasketProduct{
				ProductId: 1,
				UserId: 1,
				Number: 1,
			},
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, product *basket.BasketProduct) {
				s.EXPECT().DeleteProduct(ctx, product).Return( nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "FAIL",
			inputData: models.BasketProduct{
				ProductId: 1,
				UserId: 1,
				Number: 1,
			},
			mockBehaviorClientGRPC: func(s *mock_basket.MockBasketServiceClient,
				ctx context.Context, product *basket.BasketProduct) {
				s.EXPECT().DeleteProduct(ctx, product).Return(nil, errGetting)
			},
			expectedError: errGetting,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockClientGRPC := mock_basket.NewMockBasketServiceClient(c)
			testCase.mockBehaviorClientGRPC(mockClientGRPC, context.Background(),
				models.ModelBasketProductToGrpc(testCase.inputData))

			useCase := UseCase{
				basketServiceClient: mockClientGRPC,
			}

			assert.Equal(t, testCase.expectedError, useCase.DeleteProduct(testCase.inputData))
		})
	}
}


//func TestUseCase_PutInBasket(t *testing.T) {
//	type mockBehaviorRepository func(s *mocks.MockRepository, product models.BasketProduct)
//
//	testTable := []struct {
//		name                   string
//		inputData              models.BasketProduct
//		mockBehaviorRepository mockBehaviorRepository
//		expectedError          error
//	}{
//		{
//			name: "OK",
//			inputData: models.BasketProduct{
//				ProductId: 1,
//				Number:    2,
//			},
//			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
//				s.EXPECT().PutInBasket(product).Return(nil)
//			},
//			expectedError: nil,
//		},
//		{
//			name: "BD_ERROR",
//			inputData: models.BasketProduct{
//				ProductId: 1,
//				Number:    2,
//			},
//			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
//				s.EXPECT().PutInBasket(product).Return(errors.New(customErrors.BD_ERROR_DESCR))
//			},
//			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockRepository := mocks.NewMockRepository(c)
//
//			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)
//
//			usecase := NewBasketUseCase(mockRepository)
//
//			err := usecase.PutInBasket(testCase.inputData)
//
//			assert.Equal(t, testCase.expectedError, err)
//		})
//	}
//}
//
//func TestUseCase_GetBasket(t *testing.T) {
//	type mockBehaviorRepository func(s *mocks.MockRepository, id int)
//
//	testTable := []struct {
//		name                   string
//		inputData              int
//		mockBehaviorRepository mockBehaviorRepository
//		expectedError          error
//	}{
//		{
//			name:      "OK",
//			inputData: 1,
//			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
//				s.EXPECT().GetBasket(id).Return(nil, nil)
//			},
//			expectedError: nil,
//		},
//		{
//			name:      "BD_ERROR",
//			inputData: 2,
//			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
//				s.EXPECT().GetBasket(id).Return(nil, errors.New(customErrors.BD_ERROR_DESCR))
//			},
//			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockRepository := mocks.NewMockRepository(c)
//
//			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)
//
//			usecase := NewBasketUseCase(mockRepository)
//
//			_, err := usecase.GetBasket(testCase.inputData)
//
//			assert.Equal(t, testCase.expectedError, err)
//		})
//	}
//}
//
//func TestUseCase_DropBasket(t *testing.T) {
//	type mockBehaviorRepository func(s *mocks.MockRepository, id int)
//
//	testTable := []struct {
//		name                   string
//		inputData              int
//		mockBehaviorRepository mockBehaviorRepository
//		expectedError          error
//	}{
//		{
//			name:      "OK",
//			inputData: 1,
//			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
//				s.EXPECT().DropBasket(id).Return(nil)
//			},
//			expectedError: nil,
//		},
//		{
//			name:      "BD_ERROR",
//			inputData: 2,
//			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
//				s.EXPECT().DropBasket(id).Return(errors.New(customErrors.BD_ERROR_DESCR))
//			},
//			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockRepository := mocks.NewMockRepository(c)
//
//			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)
//
//			usecase := NewBasketUseCase(mockRepository)
//
//			err := usecase.DropBasket(testCase.inputData)
//
//			assert.Equal(t, testCase.expectedError, err)
//		})
//	}
//}
//
//func TestUseCase_DeleteProduct(t *testing.T) {
//	type mockBehaviorRepository func(s *mocks.MockRepository, product models.BasketProduct)
//
//	testTable := []struct {
//		name                   string
//		inputData              models.BasketProduct
//		mockBehaviorRepository mockBehaviorRepository
//		expectedError          error
//	}{
//		{
//			name: "OK",
//			inputData: models.BasketProduct{
//				ProductId: 1,
//				Number:    2,
//			},
//			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
//				s.EXPECT().DeleteProduct(product).Return(nil)
//			},
//			expectedError: nil,
//		},
//		{
//			name: "BD_ERROR",
//			inputData: models.BasketProduct{
//				ProductId: 1,
//				Number:    2,
//			},
//			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
//				s.EXPECT().DeleteProduct(product).Return(errors.New(customErrors.BD_ERROR_DESCR))
//			},
//			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			c := gomock.NewController(t)
//			defer c.Finish()
//
//			mockRepository := mocks.NewMockRepository(c)
//
//			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)
//
//			usecase := NewBasketUseCase(mockRepository)
//
//			err := usecase.DeleteProduct(testCase.inputData)
//
//			assert.Equal(t, testCase.expectedError, err)
//		})
//	}
//}
