package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/email"
)

type UseCase struct {
	repositoryEmail email.Repository
}

func NewEmailUseCase(repositoryEmail email.Repository) *UseCase {
	return &UseCase{
		repositoryEmail: repositoryEmail,
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

