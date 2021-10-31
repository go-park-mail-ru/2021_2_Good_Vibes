package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_order "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/order/mocks"
	mock_jwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/mocks"
	validator2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/validator"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockBehaviorPutOrder func(s *mock_order.MockUseCase, order models.Order)
type mockBehaviorParseToken func(s *mock_jwt.MockTokenManager)

func TestOrderHandler_PutOrder(t *testing.T) {
	tokenError := customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR)
	serverError := customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR)
	validError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)

	tokenErrorJson, _ := json.Marshal(tokenError)
	serverErrorJson, _ := json.Marshal(serverError)
	validErrorJson, _ := json.Marshal(validError)

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
		Date:     "28-10-2021 03:03:59",
		Address:  address,
		Cost:     50000.00,
		Status:   "new",
		Products: products,
	}

	orderJson, _ := json.Marshal(order)

	badOrder := models.Order{
		Address:  address,
		Cost:     50000.00,
		Status:   "new",
		Products: products,
	}

	badOrderJson, _ := json.Marshal(badOrder)

	tests := []struct {
		name                   string
		order                  string
		mockBehaviorPutOrder   mockBehaviorPutOrder
		mockBehaviorParseToken mockBehaviorParseToken
		expectedStatusCode     int
		expectedRequestBody    string
	}{
		{
			name:  "correct",
			order: string(orderJson),
			mockBehaviorPutOrder: func(s *mock_order.MockUseCase, order models.Order) {
				s.EXPECT().PutOrder(order).Return(1, 50000.00, nil)
			},
			mockBehaviorParseToken: func(s *mock_jwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(3), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(orderJson) + "\n",
		},
		{
			name:  "incorrect parse token",
			order: string(orderJson),
			mockBehaviorPutOrder: func(s *mock_order.MockUseCase, order models.Order) {
			},
			mockBehaviorParseToken: func(s *mock_jwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), errors.New(string(tokenErrorJson)))
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(tokenErrorJson) + "\n",
		},
		{

			name:  "incorrect put order",
			order: string(orderJson),
			mockBehaviorPutOrder: func(s *mock_order.MockUseCase, order models.Order) {
				s.EXPECT().PutOrder(order).Return(0, 50000.00, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorParseToken: func(s *mock_jwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(3), nil)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) + "\n",
		},
		{
			name:  "not valid input",
			order: string(badOrderJson),
			mockBehaviorPutOrder: func(s *mock_order.MockUseCase, order models.Order) {
			},
			mockBehaviorParseToken: func(s *mock_jwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(3), nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(validErrorJson) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newOrderUseCase := mock_order.NewMockUseCase(c)
			newJwtManager := mock_jwt.NewMockTokenManager(c)

			orderHandler := NewOrderHandler(newOrderUseCase, newJwtManager)

			tt.mockBehaviorPutOrder(newOrderUseCase, order)
			tt.mockBehaviorParseToken(newJwtManager)

			router := echo.New()
			router.POST("/cart/confirm", orderHandler.PutOrder)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/cart/confirm", bytes.NewBufferString(tt.order))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}
