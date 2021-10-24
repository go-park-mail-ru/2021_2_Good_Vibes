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
		name                string
		inputData           models.UserDataForInput
		mockBehaviorRepository mockBehaviorRepository
		mockBehaviorHasher mockBehaviorHasher
		expectedId  int
		expectedError error
	}{
		{
			name : "OK",
			inputData: models.UserDataForInput{
				Name: "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Id: 1,
					Name: "Test",
					Password: "Qwerty123.",
				}, nil)
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
				s.EXPECT().CompareHashAndPassword(hasherPassword, password)
			},
			expectedId: 1,
			expectedError: nil,
		},
		{
			name : "BD_Error",
			inputData: models.UserDataForInput{
				Name: "Test",
				Password: "Qwerty123.",
			},
			mockBehaviorRepository: func(s *mockUser.MockRepository, name string) {
				s.EXPECT().GetUserDataByName(name).Return(&models.UserDataStorage{
					Id: 1,
					Name: "Test",
					Password: "Qwerty123.",
				}, errors.New(customErrors.BD_ERROR_DESCR))
			},
			mockBehaviorHasher: func(s *mockHasher.MockHasher, hasherPassword []byte, password []byte) {
			},
			expectedId: customErrors.USER_EXISTS_ERROR,
			expectedError: errors.New(customErrors.BD_ERROR_DESCR),
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

