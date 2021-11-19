package auth

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UseCase interface {
	Login (loginUser models.UserDataForInput)(int, error)
	SignUp(signupUser models.UserDataForReg)(int, error)
}
