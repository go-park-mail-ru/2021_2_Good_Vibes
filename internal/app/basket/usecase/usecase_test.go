package usecase

import (
	"errors"
	mocks "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/mocks"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestUseCase_PutInBasket(t *testing.T) {
	type mockBehaviorRepository func(s *mocks.MockRepository, product models.BasketProduct)

	testTable := []struct {
		name                   string
		inputData              models.BasketProduct
		mockBehaviorRepository mockBehaviorRepository
		expectedError          error
	}{
		{
			name: "OK",
			inputData: models.BasketProduct{
				ProductId: 1,
				Number: 2,
			},
			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
				s.EXPECT().PutInBasket(product).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "BD_ERROR",
			inputData: models.BasketProduct{
				ProductId: 1,
				Number: 2,
			},
			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
				s.EXPECT().PutInBasket(product).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mocks.NewMockRepository(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)

			usecase := NewBasketUseCase(mockRepository)

			err := usecase.PutInBasket(testCase.inputData)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_GetBasket(t *testing.T) {
	type mockBehaviorRepository func(s *mocks.MockRepository, id int)

	testTable := []struct {
		name                   string
		inputData              int
		mockBehaviorRepository mockBehaviorRepository
		expectedError          error
	}{
		{
			name: "OK",
			inputData: 1,
			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
				s.EXPECT().GetBasket(id).Return(nil, nil)
			},
			expectedError: nil,
		},
		{
			name: "BD_ERROR",
			inputData: 2,
			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
				s.EXPECT().GetBasket(id).Return(nil, errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mocks.NewMockRepository(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)

			usecase := NewBasketUseCase(mockRepository)

			_, err := usecase.GetBasket(testCase.inputData)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_DropBasket(t *testing.T) {
	type mockBehaviorRepository func(s *mocks.MockRepository, id int)

	testTable := []struct {
		name                   string
		inputData              int
		mockBehaviorRepository mockBehaviorRepository
		expectedError          error
	}{
		{
			name: "OK",
			inputData: 1,
			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
				s.EXPECT().DropBasket(id).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "BD_ERROR",
			inputData: 2,
			mockBehaviorRepository: func(s *mocks.MockRepository, id int) {
				s.EXPECT().DropBasket(id).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mocks.NewMockRepository(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)

			usecase := NewBasketUseCase(mockRepository)

			err := usecase.DropBasket(testCase.inputData)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_DeleteProduct(t *testing.T) {
	type mockBehaviorRepository func(s *mocks.MockRepository, product models.BasketProduct)

	testTable := []struct {
		name                   string
		inputData              models.BasketProduct
		mockBehaviorRepository mockBehaviorRepository
		expectedError          error
	}{
		{
			name: "OK",
			inputData: models.BasketProduct{
				ProductId: 1,
				Number: 2,
			},
			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
				s.EXPECT().DeleteProduct(product).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "BD_ERROR",
			inputData: models.BasketProduct{
				ProductId: 1,
				Number: 2,
			},
			mockBehaviorRepository: func(s *mocks.MockRepository, product models.BasketProduct) {
				s.EXPECT().DeleteProduct(product).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mocks.NewMockRepository(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)

			usecase := NewBasketUseCase(mockRepository)

			err := usecase.DeleteProduct(testCase.inputData)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

