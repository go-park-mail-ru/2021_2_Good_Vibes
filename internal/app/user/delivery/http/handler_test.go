package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mockJwt "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/session/jwt/mocks"
	validator2 "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/validator"
	mockUser "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/mocks"
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

func TestUserHandler_SignUp(t *testing.T) {
	type mockBehaviorUseCase func(s *mockUser.MockUsecase, userReg models.UserDataForReg)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager, id int, name string)

	user1, _ := json.Marshal(models.UserDataForReg{Name: "Test1",
		Email: "test@gmail.com", Password: "Qwerty123."})
	user2, _ := json.Marshal(models.UserDataForReg{Name: "Test1",
		Email: "test@gmail.com", Password: "123"})

	user1get, _ := json.Marshal(models.UserDataProfile{Name: "Test1",
		Email: "test@gmail.com"})
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	error3get, _ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	error4get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	error5get, _ := json.Marshal(customErrors.NewError(customErrors.USER_EXISTS_ERROR, customErrors.USER_EXISTS_DESCR))
	error6get, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, customErrors.BD_ERROR_DESCR))
	error7get, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.UserDataForReg
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorSession mockBehaviorSession
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: string(user1),
			inputUser: models.UserDataForReg{
				Name:     "Test1",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
				s.EXPECT().AddUser(userReg).Return(1, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
				s.EXPECT().GetToken(id, name).Return("RandomJWT", nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(user1get) + "\n",
		},
		{
			name:      "SimplePassword",
			inputBody: string(user2),
			inputUser: models.UserDataForReg{
				Name:     "",
				Email:    "",
				Password: "",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name:      "BadJson",
			inputBody: `"username":"Test2","email":"test@gmail.com","password":"Qwerty123."}`,
			inputUser: models.UserDataForReg{
				Name:     "",
				Email:    "",
				Password: "",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error3get) + "\n",
		},
		{
			name:      "BadJsonData",
			inputBody: `{"usrname":"incorrectNameField","email":"test@gmail.com","password":"Qwerty123."}`,
			inputUser: models.UserDataForReg{
				Name:     "",
				Email:    "",
				Password: "",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error4get) + "\n",
		},
		{
			name:      "UserAlreadyExist",
			inputBody: string(user1),
			inputUser: models.UserDataForReg{
				Name:     "Test1",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
				s.EXPECT().AddUser(userReg).Return(customErrors.USER_EXISTS_ERROR, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error5get) + "\n",
		},
		{
			name:      "BDError",
			inputBody: string(user1),
			inputUser: models.UserDataForReg{
				Name:     "Test1",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
				s.EXPECT().AddUser(userReg).Return(customErrors.DB_ERROR,
					errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error6get) + "\n",
		},
		{
			name:      "GetTokenError",
			inputBody: string(user1),
			inputUser: models.UserDataForReg{
				Name:     "Test1",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForReg) {
				s.EXPECT().AddUser(userReg).Return(1, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
				s.EXPECT().GetToken(id, name).Return("RandomJWT", errors.New(customErrors.TOKEN_ERROR_DESCR))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error7get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUser := mockUser.NewMockUsecase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken, 1, testCase.inputUser.Name)
			testCase.mockBehaviorUseCase(mockUser, testCase.inputUser)

			handler := NewLoginHandler(mockUser, mockJwtToken)

			router := echo.New()
			router.POST("/signup", handler.SignUp)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/signup",
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

func TestUserHandler_Login(t *testing.T) {
	type mockBehaviorUseCase func(s *mockUser.MockUsecase, userInput models.UserDataForInput)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager, id int, name string)
	type mockBehaviorUseCaseGetUserDataByid func(s *mockUser.MockUsecase, id uint64)

	user1, _ := json.Marshal(models.UserDataForInput{Name: "Test1", Password: "Qwerty123."})

	userGet := models.UserDataProfile{Name: "Test1", Email: "test1@gmail.com"}
	user1get, _ := json.Marshal(userGet)
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.BIND_ERROR, customErrors.BIND_DESCR))
	error3get, _ := json.Marshal(customErrors.NewError(customErrors.VALIDATION_ERROR, customErrors.VALIDATION_DESCR))
	error4get, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, customErrors.BD_ERROR_DESCR))
	error5get, _ := json.Marshal(customErrors.NewError(customErrors.NO_USER_ERROR, customErrors.NO_USER_DESCR))
	error6get, _ := json.Marshal(customErrors.NewError(customErrors.WRONG_PASSWORD_ERROR, customErrors.WRONG_PASSWORD_DESCR))
	error7get, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           models.UserDataForInput
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorSession mockBehaviorSession
		mockBehaviorUseCaseGetUserDataById mockBehaviorUseCaseGetUserDataByid
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: string(user1),
			inputUser: models.UserDataForInput{
				Name:     "Test1",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
				s.EXPECT().CheckPassword(userInput).Return(1, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
				s.EXPECT().GetToken(id, name).Return("RandomJWT", nil)
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
				s.EXPECT().GetUserDataByID(id).Return(&userGet, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(user1get) + "\n",
		},
		{
			name:      "BadJson",
			inputBody: `"username":"Test2","password":"Qwerty123."}`,
			inputUser: models.UserDataForInput{
				Name:     "",
				Password: "",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userReg models.UserDataForInput) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name:      "BadJsonData",
			inputBody: `{"usrname":"incorrectNameField","password":"Qwerty123."}`,
			inputUser: models.UserDataForInput{
				Name:     "",
				Password: "",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error3get) + "\n",
		},
		{
			name:      "BD_Error",
			inputBody: string(user1),
			inputUser: models.UserDataForInput{
				Name:     "Test1",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
				s.EXPECT().CheckPassword(userInput).Return(-1, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: string(error4get) + "\n",
		},
		{
			name:      "No_user",
			inputBody: string(user1),
			inputUser: models.UserDataForInput{
				Name:     "Test1",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
				s.EXPECT().CheckPassword(userInput).Return(customErrors.NO_USER_ERROR, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error5get) + "\n",
		},
		{
			name:      "Wrong password",
			inputBody: string(user1),
			inputUser: models.UserDataForInput{
				Name:     "Test1",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
				s.EXPECT().CheckPassword(userInput).Return(customErrors.WRONG_PASSWORD_ERROR, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error6get) + "\n",
		},
		{
			name:      "GetTokenError",
			inputBody: string(user1),
			inputUser: models.UserDataForInput{
				Name:     "Test1",
				Password: "Qwerty123.",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, userInput models.UserDataForInput) {
				s.EXPECT().CheckPassword(userInput).Return(1, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager, id int, name string) {
				s.EXPECT().GetToken(id, name).Return("RandomJWT", errors.New(customErrors.TOKEN_ERROR_DESCR))
			},
			mockBehaviorUseCaseGetUserDataById: func(s *mockUser.MockUsecase, id uint64) {
				s.EXPECT().GetUserDataByID(id).Return(&userGet, nil)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: string(error7get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUser := mockUser.NewMockUsecase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken, 1, testCase.inputUser.Name)
			testCase.mockBehaviorUseCase(mockUser, testCase.inputUser)
			testCase.mockBehaviorUseCaseGetUserDataById(mockUser, uint64(1))
			handler := NewLoginHandler(mockUser, mockJwtToken)

			router := echo.New()
			router.POST("/login", handler.Login)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("POST", "/login",
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

func TestUserHandler_Profile(t *testing.T) {
	type mockBehaviorUseCase func(s *mockUser.MockUsecase, id uint64)
	type mockBehaviorSession func(s *mockJwt.MockTokenManager)

	user1get, _ := json.Marshal(models.UserDataProfile{Name: "Test", Email: "test@gmail.com"})
	error2get, _ := json.Marshal(customErrors.NewError(customErrors.TOKEN_ERROR, customErrors.TOKEN_ERROR_DESCR))
	error3get, _ := json.Marshal(customErrors.NewError(customErrors.DB_ERROR, customErrors.BD_ERROR_DESCR))

	testTable := []struct {
		name                string
		inputUser           models.UserDataProfile
		mockBehaviorUseCase mockBehaviorUseCase
		mockBehaviorSession mockBehaviorSession
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			inputUser: models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, id uint64) {
				s.EXPECT().GetUserDataByID(id).Return(&models.UserDataProfile{Name: "Test",
					Email: "test@gmail.com"}, nil)
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: string(user1get) + "\n",
		},
		{
			name: "Token Error",
			inputUser: models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, id uint64) {
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), errors.New(customErrors.TOKEN_ERROR_DESCR))
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error2get) + "\n",
		},
		{
			name: "Token Error",
			inputUser: models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			mockBehaviorUseCase: func(s *mockUser.MockUsecase, id uint64) {
				s.EXPECT().GetUserDataByID(id).Return(&models.UserDataProfile{Name: "Test",
					Email: "test@gmail.com"}, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorSession: func(s *mockJwt.MockTokenManager) {
				s.EXPECT().ParseTokenFromContext(context.Background()).Return(uint64(1), nil)
			},
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: string(error3get) + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockUser := mockUser.NewMockUsecase(c)
			mockJwtToken := mockJwt.NewMockTokenManager(c)
			testCase.mockBehaviorSession(mockJwtToken)
			testCase.mockBehaviorUseCase(mockUser, 1)

			handler := NewLoginHandler(mockUser, mockJwtToken)

			router := echo.New()
			router.GET("/profile", handler.Profile)

			val := validator.New()
			val.RegisterValidation("customPassword", validator2.Password)
			router.Validator = &validator2.CustomValidator{Validator: val}

			req := httptest.NewRequest("GET", "/profile",
				bytes.NewBufferString("empty"))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			logrus.SetOutput(ioutil.Discard)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedRequestBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockUser := mockUser.NewMockUsecase(c)
		mockJwtToken := mockJwt.NewMockTokenManager(c)

		handler := NewLoginHandler(mockUser, mockJwtToken)

		router := echo.New()
		router.GET("/logout", handler.Logout)

		val := validator.New()
		val.RegisterValidation("customPassword", validator2.Password)
		router.Validator = &validator2.CustomValidator{Validator: val}

		req := httptest.NewRequest("GET", "/logout",
			bytes.NewBufferString(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		logrus.SetOutput(ioutil.Discard)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
