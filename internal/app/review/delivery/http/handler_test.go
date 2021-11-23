package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mock_review "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/review/mocks"
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

func TestReviewHandler_AddReview(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_review.MockUseCase, review models.Review)
	type mockBehaviorParseToken func(s *mockJwt.MockTokenManager)

	tokenError := customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR)
	badJsonError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
	badJsonDataError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
	reviewAlreadyExistError := customErrors.NewError(customErrors.REVIEW_EXISTS_ERROR, customErrors.REVIEW_EXISTS_DESCR)
	serverError := customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR)

	tokenErrorJson, _ := json.Marshal(tokenError)
	badJsonErrorJson, _ := json.Marshal(badJsonError)
	badJsonDataErrorJson, _ := json.Marshal(badJsonDataError)
	reviewAlreadyExistErrorJson, _ := json.Marshal(reviewAlreadyExistError)
	serverErrorJson, _ := json.Marshal(serverError)

	newReview := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 5,
		Text: "test",
	}

	newReviewJson, _ := json.Marshal(newReview)

	testTable := []struct{
		name string
		inputBody string
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorParseToken mockBehaviorParseToken
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().AddReview(review).Return(nil)
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: string(newReviewJson) +"\n",
		},
		{
			name: "fail token ",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), errors.New("err"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedRequestBody: string(tokenErrorJson) +"\n",
		},
		{
			name: "bad json",
			inputBody: "{sdgsg",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonErrorJson) +"\n",
		},
		{
			name: "bad json data",
			inputBody: "{\"userId\": 1} ",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonDataErrorJson) +"\n",
		},
		{
			name: "review already exist",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().AddReview(review).Return(errors.New(customErrors.REVIEW_EXISTS_DESCR))
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(reviewAlreadyExistErrorJson) +"\n",
		},
		{
			name: "bd error",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().AddReview(review).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) +"\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mock_review.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			tt.mockBehaviorParseToken(mockJwtToken)
			tt.mockBehaviorUseCase(mockUseCase, newReview)

			handler := NewReviewHandler(mockUseCase, mockJwtToken)

			router := echo.New()
			router.POST("/review/add", handler.AddReview)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/review/add",
				bytes.NewBufferString(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestReviewHandler_UpdateReview(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_review.MockUseCase, review models.Review)
	type mockBehaviorParseToken func(s *mockJwt.MockTokenManager)

	tokenError := customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR)
	badJsonError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
	badJsonDataError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
	reviewNoExistError := customErrors.NewError(customErrors.NO_REVIEW_ERROR, customErrors.NO_REVIEW_DESCR)
	serverError := customErrors.NewError(customErrors.SERVER_ERROR, customErrors.BD_ERROR_DESCR)

	tokenErrorJson, _ := json.Marshal(tokenError)
	badJsonErrorJson, _ := json.Marshal(badJsonError)
	badJsonDataErrorJson, _ := json.Marshal(badJsonDataError)
	reviewNoExistErrorJson, _ := json.Marshal(reviewNoExistError)
	serverErrorJson, _ := json.Marshal(serverError)

	newReview := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 5,
		Text: "test",
	}

	newReviewJson, _ := json.Marshal(newReview)

	testTable := []struct{
		name string
		inputBody string
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorParseToken mockBehaviorParseToken
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().UpdateReview(review).Return(nil)
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: string(newReviewJson) +"\n",
		},
		{
			name: "fail token ",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), errors.New("err"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedRequestBody: string(tokenErrorJson) +"\n",
		},
		{
			name: "bad json",
			inputBody: "{sdgsg",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonErrorJson) +"\n",
		},
		{
			name: "bad json data",
			inputBody: "{\"userId\": 1} ",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonDataErrorJson) +"\n",
		},
		{
			name: "review already exist",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().UpdateReview(review).Return(errors.New(customErrors.NO_REVIEW_DESCR))
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(reviewNoExistErrorJson) +"\n",
		},
		{
			name: "bd error",
			inputBody: string(newReviewJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, review models.Review) {
				s.EXPECT().UpdateReview(review).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) +"\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mock_review.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			tt.mockBehaviorParseToken(mockJwtToken)
			tt.mockBehaviorUseCase(mockUseCase, newReview)

			handler := NewReviewHandler(mockUseCase, mockJwtToken)

			router := echo.New()
			router.POST("/review/update", handler.UpdateReview)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/review/update",
				bytes.NewBufferString(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestReviewHandler_DeleteReview(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_review.MockUseCase, userId int, productId int)
	type mockBehaviorParseToken func(s *mockJwt.MockTokenManager)

	tokenError := customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR)
	badJsonError := customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR)
	badJsonDataError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
	reviewNoExistError := customErrors.NewError(customErrors.NO_REVIEW_ERROR, customErrors.NO_REVIEW_DESCR)
	serverError := errors.New(customErrors.BD_ERROR_DESCR)

	tokenErrorJson, _ := json.Marshal(tokenError)
	badJsonErrorJson, _ := json.Marshal(badJsonError)
	badJsonDataErrorJson, _ := json.Marshal(badJsonDataError)
	reviewNoExistErrorJson, _ := json.Marshal(reviewNoExistError)
	serverErrorJson, _ := json.Marshal(serverError)

	productId := models.ProductId{
		ProductId: 1,
	}

	productIdJson, _ := json.Marshal(productId)

	testTable := []struct{
		name string
		inputBody string
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorParseToken mockBehaviorParseToken
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "ok",
			inputBody: string(productIdJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {
				s.EXPECT().DeleteReview(userId, productId).Return(nil)
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: string(productIdJson) +"\n",
		},
		{
			name: "fail token ",
			inputBody: string(productIdJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), errors.New("err"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedRequestBody: string(tokenErrorJson) +"\n",
		},
		{
			name: "bad json",
			inputBody: "{sdgsg",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonErrorJson) +"\n",
		},
		{
			name: "bad json data",
			inputBody: "{\"userId\": 1} ",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {

			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badJsonDataErrorJson) +"\n",
		},
		{
			name: "review already exist",
			inputBody: string(productIdJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {
				s.EXPECT().DeleteReview(userId, productId).Return(errors.New(customErrors.NO_REVIEW_DESCR))
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(reviewNoExistErrorJson) +"\n",
		},
		{
			name: "bd error",
			inputBody: string(productIdJson),
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userId int, productId int) {
				s.EXPECT().DeleteReview(userId, productId).Return(serverError)
			},
			mockBehaviorParseToken: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) +"\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mock_review.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			tt.mockBehaviorParseToken(mockJwtToken)
			tt.mockBehaviorUseCase(mockUseCase, 1, 1)

			handler := NewReviewHandler(mockUseCase, mockJwtToken)

			router := echo.New()
			router.POST("/review/delete", handler.DeleteReview)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/review/delete",
				bytes.NewBufferString(tt.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestReviewHandler_GetReviewsByProductId(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_review.MockUseCase, productId int)

	badQueryDataError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
	serverError := errors.New(customErrors.BD_ERROR_DESCR)

	badQueryDataErrorJson, _ := json.Marshal(badQueryDataError)
	serverErrorJson, _ := json.Marshal(serverError)

	review := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 5,
		Text: "test",
	}

	reviews := []models.Review{
		review,
	}

	reviewsJson, _ := json.Marshal(reviews)

	testTable := []struct{
		name string
		queryParam string
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "ok",
			queryParam: "1",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, productId int) {
				s.EXPECT().GetReviewsByProductId(productId).Return(reviews, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: string(reviewsJson) +"\n",
		},
		{
			name: "fail query param ",
			queryParam: "",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, productId int) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badQueryDataErrorJson) +"\n",
		},
		{
			name: "bd error",
			queryParam: "1",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, productId int) {
				s.EXPECT().GetReviewsByProductId(productId).Return(nil, serverError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) +"\n",
		},
		{
			name: "nil review",
			queryParam: "1",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, productId int) {
				s.EXPECT().GetReviewsByProductId(productId).Return(nil, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedRequestBody: string("[]") +"\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mock_review.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			tt.mockBehaviorUseCase(mockUseCase,  1)

			handler := NewReviewHandler(mockUseCase, mockJwtToken)

			router := echo.New()
			router.GET("/reviews", handler.GetReviewsByProductId)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", "/reviews?product_id=" + tt.queryParam,
				bytes.NewBufferString(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestReviewHandler_GetReviewsByUser(t *testing.T) {
	type mockBehaviorUseCase func(s *mock_review.MockUseCase, userName string)

	badQueryDataError := customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR)
	serverError := errors.New(customErrors.BD_ERROR_DESCR)

	badQueryDataErrorJson, _ := json.Marshal(badQueryDataError)
	serverErrorJson, _ := json.Marshal(serverError)

	review := models.Review{
		UserId: 1,
		UserName: "Test",
		ProductId: 1,
		Rating: 5,
		Text: "test",
	}

	reviews := []models.Review{
		review,
	}

	reviewsJson, _ := json.Marshal(reviews)

	testTable := []struct{
		name string
		queryParam string
		mockBehaviorUseCase mockBehaviorUseCase
		expectedStatusCode int
		expectedRequestBody string
	}{
		{
			name: "ok",
			queryParam: "test",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userName string) {
				s.EXPECT().GetReviewsByUser(userName).Return(reviews, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: string(reviewsJson) +"\n",
		},
		{
			name: "fail query param ",
			queryParam: "",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userName string) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedRequestBody: string(badQueryDataErrorJson) +"\n",
		},
		{
			name: "bd error",
			queryParam: "test",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userName string) {
				s.EXPECT().GetReviewsByUser(userName).Return(nil, serverError)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedRequestBody: string(serverErrorJson) +"\n",
		},
		{
			name: "nil review",
			queryParam: "test",
			mockBehaviorUseCase: func(s *mock_review.MockUseCase, userName string) {
				s.EXPECT().GetReviewsByUser(userName).Return(nil, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedRequestBody: string("[]") +"\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUseCase := mock_review.NewMockUseCase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)

			tt.mockBehaviorUseCase(mockUseCase,  "test")

			handler := NewReviewHandler(mockUseCase, mockJwtToken)

			router := echo.New()
			router.GET("/user/reviews", handler.GetReviewsByUser)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", "/user/reviews?name=" + tt.queryParam,
				bytes.NewBufferString(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedRequestBody, recorder.Body.String())
		})
	}
}