package user

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"

//go:generate mockgen -source=usecase.go -destination=mocks/usecase_mock.go
type Usecase interface {
	CheckPassword(user models.UserDataForInput) (int, error)
	GetUserDataByID(id uint64) (*models.UserDataProfile, error)
	AddUser(newUser models.UserDataForReg) (int, error)
	GenerateAvatarName() string
	SaveAvatarName(userId int, fileName string) error
	UpdateProfile(newData models.UserDataProfile) (int, error)
	UpdatePassword(newData models.UserDataPassword) error
}
