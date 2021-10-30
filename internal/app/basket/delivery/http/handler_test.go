package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	mocks "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/basket/mocks"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mockJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/mocks"
	validator2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/validator"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasketHandler_PutInBasket(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, product models.BasketProduct)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager)

	product1, _ := json.Marshal(models.BasketProduct{ProductId: 1, Number: 5})

	error2get, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))
	error3get, _ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	error4get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	error5get, _ := json.Marshal(customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		inputProduct        models.BasketProduct
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorSession mockBehaviorSession
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody:    string(product1),
			inputProduct: models.BasketProduct{ProductId: 1, Number: 5},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), nil)
			},
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.BasketProduct) {
				s.EXPECT().PutInBasket(product).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1) + "\n",
		},
		{
			name:         "Unauthorized",
			inputBody:    string(product1),
			inputProduct: models.BasketProduct{ProductId: 1, Number: 5},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), errors.New(customErrors.TOKEN_ERROR_DESCR))
			},
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.BasketProduct) {
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name:         "BadJson",
			inputBody:    "{its a bad json",
			inputProduct: models.BasketProduct{ProductId: 1, Number: 5},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), nil)
			},
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.BasketProduct) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error3get) + "\n",
		},
		{
			name:         "BadJsonData",
			inputBody:    `{"prodct_id": 0, "Numbr": 5}`,
			inputProduct: models.BasketProduct{ProductId: 0, Number: 5},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), nil)
			},
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.BasketProduct) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error4get) + "\n",
		},
		{
			name:         "BadJsonData",
			inputBody:    string(product1),
			inputProduct: models.BasketProduct{ProductId: 1, Number: 5},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), nil)
			},
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.BasketProduct) {
				s.EXPECT().PutInBasket(product).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error5get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken)
			testCase.mockBehaviorUseCase(mockUseCase, testCase.inputProduct)

			handler := NewBasketHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.POST("/cart/put", handler.PutInBasket)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/cart/put",
				bytes.NewBufferString(string(testCase.inputBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}
