package usecase

import (
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_review "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUseCase_AddReview(t *testing.T) {
	type mockBehaviorRepositoryGetReviewByUserAndProduct func(s *mock_review.MockRepository, userId int, productId int)
	type mockBehaviorRepositoryGetAllRatingOffProducts func(s *mock_review.MockRepository, productId int)
	type mockBehaviorRepositoryAddReview func(s *mock_review.MockRepository, review models.Review, rating float64)

	oldReview := models.Review{
		UserId: 0,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	review := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	productsRating := []models.ProductRating{
		{
			Rating: 1,
			Count: 1,
		},
	}
	testTable := []struct{
		name string
		inputData models.Review
		mockBehaviorRepositoryGetReviewByUserAndProduct mockBehaviorRepositoryGetReviewByUserAndProduct
		mockBehaviorRepositoryGetAllRatingOffProducts mockBehaviorRepositoryGetAllRatingOffProducts
		mockBehaviorRepositoryAddReview mockBehaviorRepositoryAddReview
		expectedError error
	}{
		{
			name: "ok",
			inputData: review,
			expectedError: nil,
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryAddReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
				s.EXPECT().AddReview(review, rating).Return(nil)
			},
		},
		{
			name: "review already exist",
			inputData: review,
			expectedError: errors.New(customErrors.REVIEW_EXISTS_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryAddReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 1",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryAddReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 2",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryAddReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 3",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryAddReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
				s.EXPECT().AddReview(review, rating).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryReviewMock := mock_review.NewMockRepository(c)
			tt.mockBehaviorRepositoryGetReviewByUserAndProduct(repositoryReviewMock,1,1)
			tt.mockBehaviorRepositoryGetAllRatingOffProducts(repositoryReviewMock, 1)
			tt.mockBehaviorRepositoryAddReview(repositoryReviewMock, review, 1)
			useCase := NewReviewUseCase(repositoryReviewMock)

			assert.Equal(t, tt.expectedError, useCase.AddReview(tt.inputData))

		})
	}
}

func TestUseCase_UpdateReview(t *testing.T) {
	type mockBehaviorRepositoryGetReviewByUserAndProduct func(s *mock_review.MockRepository, userId int, productId int)
	type mockBehaviorRepositoryGetAllRatingOffProducts func(s *mock_review.MockRepository, productId int)
	type mockBehaviorRepositoryUpdateReview func(s *mock_review.MockRepository, review models.Review, rating float64)

	oldReview := models.Review{
		UserId: 0,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	review := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	productsRating := []models.ProductRating{
		{
			Rating: 1,
			Count: 1,
		},
	}
	testTable := []struct{
		name string
		inputData models.Review
		mockBehaviorRepositoryGetReviewByUserAndProduct mockBehaviorRepositoryGetReviewByUserAndProduct
		mockBehaviorRepositoryGetAllRatingOffProducts mockBehaviorRepositoryGetAllRatingOffProducts
		mockBehaviorRepositoryUpdateReview mockBehaviorRepositoryUpdateReview
		expectedError error
	}{
		{
			name: "ok",
			inputData: review,
			expectedError: nil,
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryUpdateReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
				s.EXPECT().UpdateReview(review, rating).Return(nil)
			},
		},
		{
			name: "review not exist",
			inputData: review,
			expectedError: errors.New(customErrors.NO_REVIEW_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryUpdateReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 1",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryUpdateReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 2",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryUpdateReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
			},
		},
		{
			name: "bd error 3",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryUpdateReview: func(s *mock_review.MockRepository, review models.Review, rating float64) {
				s.EXPECT().UpdateReview(review, rating).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryReviewMock := mock_review.NewMockRepository(c)
			tt.mockBehaviorRepositoryGetReviewByUserAndProduct(repositoryReviewMock,1,1)
			tt.mockBehaviorRepositoryGetAllRatingOffProducts(repositoryReviewMock, 1)
			tt.mockBehaviorRepositoryUpdateReview(repositoryReviewMock, review, 1)
			useCase := NewReviewUseCase(repositoryReviewMock)

			assert.Equal(t, tt.expectedError, useCase.UpdateReview(tt.inputData))

		})
	}
}

func TestUseCase_DeleteReview(t *testing.T) {
	type mockBehaviorRepositoryGetReviewByUserAndProduct func(s *mock_review.MockRepository, userId int, productId int)
	type mockBehaviorRepositoryGetAllRatingOffProducts func(s *mock_review.MockRepository, productId int)
	type mockBehaviorRepositoryDeleteReview func(s *mock_review.MockRepository, userId int, productId int, rating float64)

	oldReview := models.Review{
		UserId: 0,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	review := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 1,
		Text: "test",
	}

	productsRating := []models.ProductRating{
		{
			Rating: 1,
			Count: 1,
		},
	}
	testTable := []struct{
		name string
		inputData models.Review
		mockBehaviorRepositoryGetReviewByUserAndProduct mockBehaviorRepositoryGetReviewByUserAndProduct
		mockBehaviorRepositoryGetAllRatingOffProducts mockBehaviorRepositoryGetAllRatingOffProducts
		mockBehaviorRepositoryDeleteReview mockBehaviorRepositoryDeleteReview
		expectedError error
	}{
		{
			name: "ok",
			inputData: review,
			expectedError: nil,
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryDeleteReview: func(s *mock_review.MockRepository, userId int, productId int, rating float64) {
				s.EXPECT().DeleteReview(userId, productId, float64(0)).Return(nil)
			},
		},
		{
			name: "review not exist",
			inputData: review,
			expectedError: errors.New(customErrors.NO_REVIEW_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryDeleteReview: func(s *mock_review.MockRepository, userId int, productId int, rating float64) {
			},
		},
		{
			name: "bd error 1",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(oldReview, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
			},
			mockBehaviorRepositoryDeleteReview: func(s *mock_review.MockRepository, userId int, productId int, rating float64) {
			},
		},
		{
			name: "bd error 2",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorRepositoryDeleteReview: func(s *mock_review.MockRepository, userId int, productId int, rating float64) {
			},
		},
		{
			name: "bd error 3",
			inputData: review,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
			mockBehaviorRepositoryGetReviewByUserAndProduct: func(s *mock_review.MockRepository, userId int, productId int) {
				s.EXPECT().GetReviewByUserAndProduct(userId, productId).Return(review, nil)
			},
			mockBehaviorRepositoryGetAllRatingOffProducts: func(s *mock_review.MockRepository, productId int) {
				s.EXPECT().GetAllRatingsOfProduct(productId).Return(productsRating, nil)
			},
			mockBehaviorRepositoryDeleteReview: func(s *mock_review.MockRepository, userId int, productId int, rating float64) {
				s.EXPECT().DeleteReview(userId, productId, float64(0)).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repositoryReviewMock := mock_review.NewMockRepository(c)
			tt.mockBehaviorRepositoryGetReviewByUserAndProduct(repositoryReviewMock,1,1)
			tt.mockBehaviorRepositoryGetAllRatingOffProducts(repositoryReviewMock, 1)
			tt.mockBehaviorRepositoryDeleteReview(repositoryReviewMock,  1, 1, float64(1))
			useCase := NewReviewUseCase(repositoryReviewMock)

			assert.Equal(t, tt.expectedError, useCase.DeleteReview(1, 1))

		})
	}
}