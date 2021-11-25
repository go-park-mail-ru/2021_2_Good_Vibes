package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mocks "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/mocks"
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

func TestProductHandler_AddProduct(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, product models.Product)

	product1 := models.Product{
		Id:           1,
		Image:        "cartinka",
		Name:         "cartinka",
		Price:        2280,
		Rating:       6,
		CountInStock: 50,
		Description:  "OPICANIE",
	}
	product1Send, _ := json.Marshal(product1)
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	error3get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	error4get, _ := json.Marshal(customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		inputProduct        models.Product
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody:    string(product1Send),
			inputProduct: product1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.Product) {
				s.EXPECT().AddProduct(product).Return(1, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "BadJson",
			inputBody:    "{bad json}",
			inputProduct: product1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.Product) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name:         "BadJsonData",
			inputBody:    `{"id": 123,"imggg": "sdfsg"}`,
			inputProduct: product1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.Product) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error3get) + "\n",
		},
		{
			name:         "BDERROR",
			inputBody:    string(product1Send),
			inputProduct: product1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, product models.Product) {
				s.EXPECT().AddProduct(product).Return(1, errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error4get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)
			testCase.mockBehaviorUseCase(mockUseCase, testCase.inputProduct)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.POST("/product/add", handler.AddProduct)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/product/add",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestProductHandler_GetProductById(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, id int)

	product1 := models.Product{
		Id:           1,
		Image:        "cartinka",
		Name:         "cartinka",
		Price:        2280,
		Rating:       6,
		CountInStock: 50,
		Description:  "OPICANIE",
	}
	product1Send, _ := json.Marshal(product1)
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		target              string
		inputProduct        int
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody:    string(product1Send),
			target:       "/product?id=1",
			inputProduct: 1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, id int) {
				s.EXPECT().GetProductById(id).Return(product1, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "BadQueryParamName",
			inputBody:    string(product1Send),
			target:       "/product?ids=1",
			inputProduct: 1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, id int) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name:         "BadQueryParamValues",
			inputBody:    string(product1Send),
			target:       "/product?id=adf",
			inputProduct: 1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase, id int) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error2get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)
			testCase.mockBehaviorUseCase(mockUseCase, testCase.inputProduct)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.GET("/product", handler.GetProductById)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", testCase.target,
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase)

	product1 := models.Product{
		Id:           1,
		Image:        "cartinka",
		Name:         "cartinka",
		Price:        2280,
		Rating:       6,
		CountInStock: 50,
		Description:  "OPICANIE",
	}
	products := []models.Product{product1}
	product1Send, _ := json.Marshal(products)
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		inputProduct        int
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody:    string(product1Send),
			inputProduct: 1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase) {
				s.EXPECT().GetAllProducts().Return(products, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "BDERROR",
			inputBody:    string(product1Send),
			inputProduct: 1,
			mockBehaviorUseCase: func(s *mocks.MockUseCase) {
				s.EXPECT().GetAllProducts().Return(products, errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error2get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)
			testCase.mockBehaviorUseCase(mockUseCase)

			mockJwtToken := mockJwt.NewMockTokenManager(c)
			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.GET("/homepage", handler.GetAllProducts)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", "/homepage",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestProductHandler_GetFavouriteProducts(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, userId int)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager)

	product1 := models.Product{
		Id:           1,
		Image:        "cartinka",
		Name:         "cartinka",
		Price:        2280,
		Rating:       6,
		CountInStock: 50,
		Description:  "OPICANIE",
	}
	products := []models.Product{product1}
	product1Send, _ := json.Marshal(products)
	errorTokenGet := errors.New(customErrors.TOKEN_ERROR_DESCR)
	errorTokenGetJson, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, errorTokenGet.Error()))
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		mockBehaviorSession mockBehaviorSession
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().GetFavouriteProducts(userId).Return(products, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "parse token error",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), errorTokenGet)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(errorTokenGetJson) + "\n",
		},
		{
			name:         "bd error ",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().GetFavouriteProducts(userId).Return(products, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error2get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)

			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken)
			testCase.mockBehaviorUseCase(mockUseCase, 1)
			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.GET("/product/favorite/get", handler.GetFavouriteProducts)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", "/product/favorite/get",
				bytes.NewBufferString(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestProductHandler_DeleteFavouriteProduct(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, userId int)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager)

	product1 := models.FavouriteProduct{
		Id:           1,
		UserId: 1,
	}
	product1Send, _ := json.Marshal(product1)
	errorTokenGet := errors.New(customErrors.TOKEN_ERROR_DESCR)
	errorTokenGetJson, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, errorTokenGet.Error()))
	errorBindJson,_ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	errorValidationJson,_ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	errorBdError, _ := json.Marshal(customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody 			string
		mockBehaviorSession mockBehaviorSession
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().DeleteFavouriteProduct(product1).Return( nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "parse token error",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), errorTokenGet)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(errorTokenGetJson) + "\n",
		},
		{
			name:         "bad json",
			inputBody: "{rsfdhgf",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errorBindJson) + "\n",
		},
		{
			name:         "validate error",
			inputBody: "{\"rsfdhgf\": 1}",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errorValidationJson) + "\n",
		},
		{
			name:         "bd error ",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().DeleteFavouriteProduct(product1).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(errorBdError) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)

			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken)
			testCase.mockBehaviorUseCase(mockUseCase, 1)
			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.POST("/product/favorite/delete", handler.DeleteFavouriteProduct)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/product/favorite/delete",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestProductHandler_AddFavouriteProduct(t *testing.T) {
	type mockBehaviorUseCase func(s *mocks.MockUseCase, userId int)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager)

	product1 := models.FavouriteProduct{
		Id:           1,
		UserId: 1,
	}
	product1Send, _ := json.Marshal(product1)
	errorTokenGet := errors.New(customErrors.TOKEN_ERROR_DESCR)
	errorTokenGetJson, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, errorTokenGet.Error()))
	errorBindJson,_ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	errorValidationJson,_ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	errorBdError, _ := json.Marshal(customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody 			string
		mockBehaviorSession mockBehaviorSession
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().AddFavouriteProduct(product1).Return( nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(product1Send) + "\n",
		},
		{
			name:         "parse token error",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(0), errorTokenGet)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(errorTokenGetJson) + "\n",
		},
		{
			name:         "bad json",
			inputBody: "{rsfdhgf",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errorBindJson) + "\n",
		},
		{
			name:         "validate error",
			inputBody: "{\"rsfdhgf\": 1}",
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errorValidationJson) + "\n",
		},
		{
			name:         "bd error ",
			inputBody: string(product1Send),
			mockBehaviorUseCase: func(s *mocks.MockUseCase, userId int) {
				s.EXPECT().AddFavouriteProduct(product1).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(errorBdError) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mocks.NewMockUseCase(c)

			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken)
			testCase.mockBehaviorUseCase(mockUseCase, 1)
			handler := NewProductHandler(mockUseCase, mockJwtToken)
			router := echo.New()
			router.POST("/product/favorite/add", handler.AddFavouriteProduct)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/product/favorite/add",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

