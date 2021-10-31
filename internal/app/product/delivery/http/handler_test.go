package http

import (
	"bytes"
	"encoding/json"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mocks "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/mocks"
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

			handler := NewProductHandler(mockUseCase)
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

			handler := NewProductHandler(mockUseCase)
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

			handler := NewProductHandler(mockUseCase)
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
