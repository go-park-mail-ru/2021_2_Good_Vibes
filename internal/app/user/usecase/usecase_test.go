package usecase

import (
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	mockHasher "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher/mock"
	mockUser "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestUseCase_CheckPassword(t *testing.T) {
	type mockBehaviorRepository func(s *mockUser.MockRepository, name string)
	type mockBehaviorHasher func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte)

	testTable := []struct {
		name                   string
		inputData              models.UserDataForInput
		mockBehaviorRepository mockBehaviorRepository
		mockBehaviorHasher     mockBehaviorHasher
		expectedId             int
		expectedError          error
	}{
		{
			name: "OK",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Id:       1,
					Name:     "Test",
					Password: "Qwerty123.",
				}, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
				s.EXPECT().CompareHashAndPassword(hasherPassword, password)
			},
			expectedId:    1,
			expectedError: nil,
		},
		{
			name: "BD_Error",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Id:       1,
					Name:     "Test",
					Password: "Qwerty123.",
				}, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
			},
			expectedId:    customErrors.USER_EXISTS_ERROR,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
		{
			name: "No user",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(nil, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
			},
			expectedId:    customErrors.NO_USER_ERROR,
			expectedError: nil,
		},
		{
			name: "Wrong Password",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Id:       1,
					Name:     "Test",
					Password: "Qwerty123.",
				}, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
				s.EXPECT().CompareHashAndPassword(hasherPassword, password).Return(errors.New(customErrors.WRONG_PASSWORD_DESCR))
			},
			expectedId:    customErrors.WRONG_PASSWORD_ERROR,
			expectedError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mockUser.NewMockRepository(c)
			mockHasher := mockHasher.NewMockHasher(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData.Name)
			testCase.mockBehaviorHasher(mockHasher, []byte(testCase.inputData.Password), []byte(testCase.inputData.Password))

			usecase := NewUsecase(mockRepository, mockHasher)

			id, err := usecase.CheckPassword(testCase.inputData)

			assert.Equal(t, testCase.expectedId, id)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_AddUser(t *testing.T) {
	type mockBehaviorRepositoryGetUserDataByName func(s *mockUser.MockRepository, name string)
	type mockBehaviorRepositoryInsertUser func(s *mockUser.MockRepository, user models.UserDataForReg)
	type mockBehaviorHasher func(s *mockHasher.MockHasher, password []byte)

	testTable := []struct {
		name                                    string
		inputData                               models.UserDataForInput
		inputDataForReg                         models.UserDataForReg
		mockBehaviorRepositoryGetUserDataByName mockBehaviorRepositoryGetUserDataByName
		mockBehaviorRepositoryInsertUser        mockBehaviorRepositoryInsertUser
		mockBehaviorHasher                      mockBehaviorHasher
		expectedId                              int
		expectedError                           error
	}{
		{
			name: "OK",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			inputDataForReg: models.UserDataForReg{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorRepositoryGetUserDataByName: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(nil, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, password []byte) {
				s.EXPECT().GenerateFromPassword(password).Return(password, nil)
			},
			mockBehaviorRepositoryInsertUser: func(s *mockUser.MockRepository, user models.UserDataForReg) {
				s.EXPECT().InsertUser(user).Return(1, nil)
			},
			expectedId:    1,
			expectedError: nil,
		},
		{
			name: "DB_error_GetUser",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			inputDataForReg: models.UserDataForReg{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorRepositoryGetUserDataByName: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(nil, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, password []byte) {
			},
			mockBehaviorRepositoryInsertUser: func(s *mockUser.MockRepository, user models.UserDataForReg) {
			},
			expectedId:    customErrors.SERVER_ERROR,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
		{
			name: "User Exist",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			inputDataForReg: models.UserDataForReg{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorRepositoryGetUserDataByName: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Name: "Ya ushe est'",
				}, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, password []byte) {
			},
			mockBehaviorRepositoryInsertUser: func(s *mockUser.MockRepository, user models.UserDataForReg) {
			},
			expectedId:    customErrors.USER_EXISTS_ERROR,
			expectedError: nil,
		},
		{
			name: "HasherError",
			inputData: models.UserDataForInput{
				Name:     "Test",
				Password: "Qwerty123.",
			},
			inputDataForReg: models.UserDataForReg{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "Qwerty123.",
			},
			mockBehaviorRepositoryGetUserDataByName: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(nil, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, password []byte) {
				s.EXPECT().GenerateFromPassword(password).Return(password, errors.New(customErrors.HASHER_ERROR_DESCR))
			},
			mockBehaviorRepositoryInsertUser: func(s *mockUser.MockRepository, user models.UserDataForReg) {
			},
			expectedId:    customErrors.SERVER_ERROR,
			expectedError: errors.New(customErrors.HASHER_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mockUser.NewMockRepository(c)
			mockHasher := mockHasher.NewMockHasher(c)

			testCase.mockBehaviorRepositoryGetUserDataByName(mockRepository, testCase.inputData.Name)
			testCase.mockBehaviorHasher(mockHasher, []byte(testCase.inputData.Password))
			testCase.mockBehaviorRepositoryInsertUser(mockRepository, testCase.inputDataForReg)

			usecase := NewUsecase(mockRepository, mockHasher)

			id, err := usecase.AddUser(testCase.inputDataForReg)

			assert.Equal(t, testCase.expectedId, id)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_GetUserDataByID(t *testing.T) {
	type mockBehaviorRepository func(s *mockUser.MockRepository, id uint64)

	testTable := []struct {
		name                   string
		inputData              uint64
		outputData             models.UserDataProfile
		mockBehaviorRepository mockBehaviorRepository
		expected               *models.UserDataProfile
		expectedError          error
	}{
		{
			name:      "OK",
			inputData: uint64(1),
			outputData: models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, id uint64) {
				s.EXPECT().GetUserDataById(id).Return(&models.UserDataStorage{
					Name:  "Test",
					Email: "test@gmail.com",
				}, nil)
			},
			expected: &models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			expectedError: nil,
		},
		{
			name:      "DB_error",
			inputData: uint64(1),
			outputData: models.UserDataProfile{
				Name:  "Test",
				Email: "test@gmail.com",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, id uint64) {
				s.EXPECT().GetUserDataById(id).Return(&models.UserDataStorage{
					Name:  "Test",
					Email: "test@gmail.com",
				}, errors.New(customErrors.BD_ERROR_DESCR))
			},
			expected:      nil,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mockUser.NewMockRepository(c)
			mockHasher := mockHasher.NewMockHasher(c)

			testCase.mockBehaviorRepository(mockRepository, testCase.inputData)

			usecase := NewUsecase(mockRepository, mockHasher)

			user, err := usecase.GetUserDataByID(testCase.inputData)

			assert.Equal(t, testCase.expected, user)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUseCase_SaveAvatarName(t *testing.T) {
	type mockBehaviorRepository func(s *mockUser.MockRepository, userId int, fileName string)

	testTable := []struct {
		name      string
		inputData struct {
			userId   int
			fileName string
		}
		mockBehaviorRepository mockBehaviorRepository
		expectedError          error
	}{
		{
			name: "OK",
			inputData: struct {
				userId   int
				fileName string
			}{userId: 1, fileName: "test"},
			mockBehaviorRepository: func(s *mockUser.MockRepository, userId int, fileName string) {
				s.EXPECT().SaveAvatarName(userId, fileName).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "DB_error",
			inputData: struct {
				userId   int
				fileName string
			}{userId: 1, fileName: "test"},
			mockBehaviorRepository: func(s *mockUser.MockRepository, userId int, fileName string) {
				s.EXPECT().SaveAvatarName(userId, fileName).Return(errors.New(customErrors.BD_ERROR_DESCR))
			},
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mockUser.NewMockRepository(c)
			mockHasher := mockHasher.NewMockHasher(c)

			testCase.mockBehaviorRepository(mockRepository,
				testCase.inputData.userId, testCase.inputData.fileName)

			usecase := NewUsecase(mockRepository, mockHasher)

			err := usecase.SaveAvatarName(testCase.inputData.userId, testCase.inputData.fileName)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
