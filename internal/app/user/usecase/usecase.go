package usecase

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	guuid "github.com/google/uuid"
)

type usecase struct {
	repository user.Repository
	hasher     hasher.Hasher
}

func NewUsecase(repositoryUser user.Repository, hasher hasher.Hasher) *usecase {
	return &usecase{
		repository: repositoryUser,
		hasher:     hasher,
	}
}

func (us *usecase) CheckPassword(user models.UserDataForInput) (int, error) {
	userFromDb, err := us.repository.GetUserDataByName(user.Name)
	if err != nil {
		return customErrors.USER_EXISTS_ERROR, err
	}

	if userFromDb == nil {
		return customErrors.NO_USER_ERROR, nil
	}

	if err = us.hasher.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password)); err != nil {
		return customErrors.WRONG_PASSWORD_ERROR, nil
	}

	return userFromDb.Id, nil
}

func (us *usecase) AddUser(newUser models.UserDataForReg) (int, error) {
	id, err := us.repository.GetUserDataByName(newUser.Name)
	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	if id != nil {
		return customErrors.USER_EXISTS_ERROR, nil
	}

	passwordHash, err := us.hasher.GenerateFromPassword([]byte(newUser.Password))
	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	newUser.Password = string(passwordHash)

	return us.repository.InsertUser(newUser)
}

func (us *usecase) GetUserDataByID(id uint64) (*models.UserDataProfile, error) {
	userDataStorage, err := us.repository.GetUserDataById(id)
	if err != nil {
		return nil, err
	}

	var userProfile models.UserDataProfile
	userProfile.Name = userDataStorage.Name
	userProfile.Email = userDataStorage.Email
	return &userProfile, nil
}

func (us *usecase) GenerateAvatarName() string {
	return guuid.New().String()
}

func (us *usecase) SaveAvatarName(userId int, fileName string) error {
	err := us.repository.SaveAvatarName(userId, fileName)
	if err != nil {
		return err
	}

	return nil
}
