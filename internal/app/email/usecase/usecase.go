package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/email"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	guuid "github.com/google/uuid"
)

type UseCase struct {
	repositoryEmail email.Repository
	repositoryUser user.Repository
}

func NewEmailUseCase(repositoryEmail email.Repository, repositoryUser user.Repository) *UseCase {
	return &UseCase{
		repositoryEmail: repositoryEmail,
		repositoryUser: repositoryUser,
	}
}

func (uc *UseCase) ConfirmEmail(email string, token string) error {
	repoToken, err := uc.repositoryEmail.GetToken(email)
	if err != nil {
		return err
	}

	if repoToken != token {
		err := errors.New("Invalid token, ")
		return err
	}

	err = uc.repositoryEmail.ConfirmEmail(email)
	if err != nil {
		return err
	}
	return  nil
}

func (uc *UseCase) SendConfirmationEmail(email string, token string) error {
	repoToken, err := uc.repositoryEmail.GetToken(email)
	if err != nil {
		return err
	}

	if repoToken != token {
		err := errors.New("Invalid token, ")
		return err
	}

	err = uc.repositoryEmail.ConfirmEmail(email)
	if err != nil {
		return err
	}
	return  nil
}

func (uc *UseCase) GetUserEmailById(id uint64) (string, error) {
	userData, err := uc.repositoryUser.GetUserDataById(id)
	if err != nil {
		return "", err
	}

	return  userData.Email, nil
}

func (uc *UseCase) GenerateToken() string {
	return guuid.New().String()
}

func (uc *UseCase) InsertUserToken(email string, token string) error {
	err := uc.repositoryUser.InsertUserToken(email, token)
	if err != nil {
		return err
	}

	return nil
}
