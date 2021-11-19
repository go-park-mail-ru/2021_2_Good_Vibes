package usecase

import (
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/hasher"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/microservice/auth"
)

type usecase struct {
	repository auth.Repository
	hasher hasher.Hasher
}

func NewUsecase(repositoryUser auth.Repository, hasher hasher.Hasher) *usecase {
	return &usecase{
		repository: repositoryUser,
		hasher:     hasher,
	}
}

func (us *usecase) Login(user models.UserDataForInput) (int, error) {
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


func (us *usecase) SignUp(newUser models.UserDataForReg) (int, error) {
	id, err := us.repository.GetUserDataByName(newUser.Name)
	if err != nil {
		return customErrors.DB_ERROR, errors.New("db error")
	}

	if id != nil {
		return customErrors.USER_EXISTS_ERROR, nil
	}

	passwordHash, err := us.hasher.GenerateFromPassword([]byte(newUser.Password))
	if err != nil {
		return customErrors.SERVER_ERROR, errors.New("server error")
	}

	newUser.Password = string(passwordHash)

	return us.repository.InsertUser(newUser)
}