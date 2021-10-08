package user

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

type UserUseCase interface {
	CheckPassword(user models.UserDataForInput) (int, error)
	AddUser(newUser models.UserDataForReg) (int, error)
}
