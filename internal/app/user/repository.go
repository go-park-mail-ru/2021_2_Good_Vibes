package user

import (
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type Repository interface {
	GetUserDataByName(name string) (*models.UserDataStorage, error)
	InsertUser(newUser models.UserDataForReg) (int, error)
}
