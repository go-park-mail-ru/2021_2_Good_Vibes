package http

import (
	"bytes"
	"encoding/json"
	mock_category "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/mocks"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	validator2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/validator"
	"github.com/go-playground/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderHandler_GetCategoryProducts(t *testing.T) {
	type mockBehaviorGetProductsByCategory func(s *mock_category.MockUseCase)

	err := customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR)
	errJson, _ := json.Marshal(err)

	defaultFilter := postgre.Filter{
		NameCategory: "clothes",
		MaxPrice:     1000000,
		MaxRating:    5,
		OrderBy:      "rating",
		TypeOrder:    "desc",
	}
	products := []models.Product{
		{
			Id:           1,
			Image:        "image",
			Name:         "name",
			Price:        1000.0,
			Rating:       5,
			Category:     "CLOTHES",
			CountInStock: 1000,
			Description:  "cool description",
		},
	}

	productsJson, _ := json.Marshal(products)

	tests := []struct {
		name                              string
		mockBehaviorGetProductsByCategory mockBehaviorGetProductsByCategory
		expectedStatusCode                int
		expectedRequestBody               string
	}{
		{
			name: "correct",
			mockBehaviorGetProductsByCategory: func(s *mock_category.MockUseCase) {
				s.EXPECT().GetProductsByCategory(defaultFilter).Return(products, nil)
			},

			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(productsJson) + "\n",
		},
		{
			name: "error",
			mockBehaviorGetProductsByCategory: func(s *mock_category.MockUseCase) {
				s.EXPECT().GetProductsByCategory(defaultFilter).Return(nil, errors.New(customErrors.BD_ERROR_DESCR))
			},

			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errJson) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newCategoryUseCase := mock_category.NewMockUseCase(c)

			categoryHandler := NewCategoryHandler(newCategoryUseCase)

			tt.mockBehaviorGetProductsByCategory(newCategoryUseCase)

			router := echo.New()
			router.GET("/category/:name", categoryHandler.GetCategoryProducts)

			req := httptest.NewRequest("GET", "/category/clothes", nil)

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	type mockBehaviorCreateCategory func(s *mock_category.MockUseCase)
	type mockBehaviorGetAllCategories func(s *mock_category.MockUseCase)

	categories := models.CategoryNode{
		Name:     "CLOTHES",
		Nesting:  0,
		Children: nil,
	}

	categoriesJson, _ := json.Marshal(categories)

	newCategory := models.CreateCategory{
		Category:       "MEN_CLOTHES",
		ParentCategory: "CLOTHES",
	}

	newCategoryJson, _ := json.Marshal(newCategory)

	err := customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR)
	errJson, _ := json.Marshal(err)

	tests := []struct {
		name                         string
		mockBehaviorCreateCategory   mockBehaviorCreateCategory
		mockBehaviorGetAllCategories mockBehaviorGetAllCategories
		expectedStatusCode           int
		expectedRequestBody          string
	}{
		{
			name: "correct",
			mockBehaviorCreateCategory: func(s *mock_category.MockUseCase) {
				s.EXPECT().CreateCategory(newCategory).Return(nil)
			},
			mockBehaviorGetAllCategories: func(s *mock_category.MockUseCase) {
				s.EXPECT().GetAllCategories().Return(categories, nil)
			},

			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(categoriesJson) + "\n",
		},
		{
			name: "error create",
			mockBehaviorCreateCategory: func(s *mock_category.MockUseCase) {
				s.EXPECT().CreateCategory(newCategory).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorGetAllCategories: func(s *mock_category.MockUseCase) {
			},

			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(errJson) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newCategoryUseCase := mock_category.NewMockUseCase(c)

			categoryHandler := NewCategoryHandler(newCategoryUseCase)

			tt.mockBehaviorCreateCategory(newCategoryUseCase)
			tt.mockBehaviorGetAllCategories(newCategoryUseCase)

			router := echo.New()
			router.POST("/category/create", categoryHandler.CreateCategory)

			req := httptest.NewRequest("POST", "/category/create", bytes.NewBufferString(string(newCategoryJson)))

			val := validator.New()
			router.Validator = &validator2.CustomValidator{Validator: val}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestCategoryHandler_GetCategories(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_category.MockUseCase)

	categories := models.CategoryNode{
		Name:     "CLOTHES",
		Nesting:  0,
		Children: nil,
	}
	error1 := errors.New("bdError")
	categoriesJson, _ := json.Marshal(categories)
	error1Json, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, error1.Error()))

	tests := []struct {
		name                string
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode  int
		expectedStatusBody  string
	}{
		{
			name: "ok",
			mockBehaviorUseCase: func(s *mock_category.MockUseCase) {
				s.EXPECT().GetAllCategories().Return(categories, nil)
			},
			expectedStatusCode: 200,
			expectedStatusBody: string(categoriesJson) + "\n",
		},
		{
			name: "bd_error",
			mockBehaviorUseCase: func(s *mock_category.MockUseCase) {
				s.EXPECT().GetAllCategories().Return(categories, error1)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedStatusBody: string(error1Json) + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newCategoryUseCase := mock_category.NewMockUseCase(c)
			tt.mockBehaviorUseCase(newCategoryUseCase)
			categoryHandler := NewCategoryHandler(newCategoryUseCase)

			router := echo.New()
			router.GET("/category", categoryHandler.GetCategories)

			req := httptest.NewRequest("GET", "/category", bytes.NewBufferString(""))

			val := validator.New()
			router.Validator = &validator2.CustomValidator{Validator: val}

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedStatusBody, recorder.Body.String())
		})
	}
}
