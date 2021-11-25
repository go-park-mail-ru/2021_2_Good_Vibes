package usecase

import (
	"errors"
	mock_category "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/category/mocks"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_product "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/mocks"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_CreateCategory(t *testing.T) {
	type mockBehaviorRepositoryCategory func(s *mock_category.MockRepository, category models.CreateCategory)

	newCategory := models.CreateCategory{
		Category:       "MEN_CLOTHES",
		ParentCategory: "CLOTHES",
	}

	errorFail := errors.New("bd Error")
	tests := []struct {
		name                           string
		mockBehaviorRepositoryCategory mockBehaviorRepositoryCategory
		expectedError                  error
	}{
		{
			name: "ok",
			mockBehaviorRepositoryCategory: func(s *mock_category.MockRepository, category models.CreateCategory) {
				s.EXPECT().CreateCategory(category).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "fail",
			mockBehaviorRepositoryCategory: func(s *mock_category.MockRepository, category models.CreateCategory) {
				s.EXPECT().CreateCategory(category).Return(errorFail)
			},
			expectedError: errorFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryCategoryMock := mock_category.NewMockRepository(c)
			repositoryProductMock := mock_product.NewMockRepository(c)
			tt.mockBehaviorRepositoryCategory(repositoryCategoryMock, newCategory)
			useCase := NewCategoryUseCase(repositoryCategoryMock, repositoryProductMock)

			assert.Equal(t, tt.expectedError, useCase.CreateCategory(newCategory))
		})
	}
}

func TestUseCase_GetAllCategories(t *testing.T) {
	type mockBehaviorRepositoryCategory func(s *mock_category.MockRepository)

	categories := models.CategoryNode{
		Name:     "CLOTHES",
		Nesting:  0,
		Children: nil,
	}

	nestingCategories := []models.NestingCategory{
		{
			Nesting:     0,
			Name:        "CLOTHES",
			Description: "",
		},
	}
	errorFail := errors.New("bd Error")
	tests := []struct {
		name                           string
		mockBehaviorRepositoryCategory mockBehaviorRepositoryCategory
		expectedData                   models.CategoryNode
		expectedError                  error
	}{
		{
			name: "ok",
			mockBehaviorRepositoryCategory: func(s *mock_category.MockRepository) {
				s.EXPECT().SelectAllCategories().Return(nestingCategories, nil)
			},
			expectedError: nil,
			expectedData:  categories,
		},
		{
			name: "fail",
			mockBehaviorRepositoryCategory: func(s *mock_category.MockRepository) {
				s.EXPECT().SelectAllCategories().Return(nestingCategories, errorFail)
			},
			expectedError: errorFail,
			expectedData:  models.CategoryNode{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryCategoryMock := mock_category.NewMockRepository(c)
			repositoryProductMock := mock_product.NewMockRepository(c)
			tt.mockBehaviorRepositoryCategory(repositoryCategoryMock)
			useCase := NewCategoryUseCase(repositoryCategoryMock, repositoryProductMock)

			data, error := useCase.GetAllCategories()

			assert.Equal(t, tt.expectedError, error)
			assert.Equal(t, tt.expectedData, data)
		})
	}
}

func TestUseCase_GetProductsByCategory(t *testing.T) {
	type mockBehaviorRepositoryProduct func(s *mock_product.MockRepository, filter postgre.Filter)

	defaultFilter := postgre.Filter{
		NameCategory: "clothes",
		MaxPrice:     1000000,
		MaxRating:    5,
		OrderBy:      "rating",
		TypeOrder:    "desc",
	}

	products := []models.Product{
		{
			Id:    0,
			Image: "",
			Name:  "test",
		},
	}
	errorFail := errors.New("bd Error")
	tests := []struct {
		name                          string
		mockBehaviorRepositoryProduct mockBehaviorRepositoryProduct
		expectedData                  []models.Product
		expectedError                 error
	}{
		{
			name: "ok",
			mockBehaviorRepositoryProduct: func(s *mock_product.MockRepository, filter postgre.Filter) {
				s.EXPECT().GetByCategory(defaultFilter).Return(products, nil)
			},
			expectedError: nil,
			expectedData:  products,
		},
		{
			name: "fail",
			mockBehaviorRepositoryProduct: func(s *mock_product.MockRepository, filter postgre.Filter) {
				s.EXPECT().GetByCategory(defaultFilter).Return(nil, errorFail)
			},
			expectedError: errorFail,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryCategoryMock := mock_category.NewMockRepository(c)
			repositoryProductMock := mock_product.NewMockRepository(c)
			tt.mockBehaviorRepositoryProduct(repositoryProductMock, defaultFilter)
			useCase := NewCategoryUseCase(repositoryCategoryMock, repositoryProductMock)

			data, error := useCase.GetProductsByCategory(defaultFilter)

			assert.Equal(t, tt.expectedError, error)
			assert.Equal(t, tt.expectedData, data)
		})
	}
}
