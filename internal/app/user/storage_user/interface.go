package storage_user

import "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"

type UserUseCase interface {
	IsUserExists(user user.UserInput) (int, error)
	AddUser(newUser user.User) (int, error)
}
