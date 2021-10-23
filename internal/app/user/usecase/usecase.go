package usecase

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type usecase struct {
	repository user.Repository
}

func NewUsecase(repositoryUser user.Repository) *usecase {
	return &usecase{
		repository: repositoryUser,
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

	if err = bcrypt.CompareHashAndPassword([]byte(userFromDb.Password), []byte(user.Password)); err != nil {
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	newUser.Password = string(passwordHash)

	return us.repository.InsertUser(newUser)
}

func (us *usecase) GetUserDataByID(id uint64) (*models.UserDataStorage, error) {
	return us.repository.GetUserDataById(id)
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
