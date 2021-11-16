package user

import (
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go
type Repository interface {
	GetUserDataByName(name string) (*models.UserDataStorage, error)
	GetUserDataById(id uint64) (*models.UserDataStorage, error)
	InsertUser(newUser models.UserDataForReg) (int, error)
	SaveAvatarName(userId int, fileName string) error
	UpdateUser (newData models.UserDataProfile) error
	UpdatePassword(newData models.UserDataPassword) error
}
