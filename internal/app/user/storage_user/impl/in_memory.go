package impl

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	userModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type StorageUserMemory struct {
	mx1      sync.RWMutex
	mx2      sync.RWMutex
	storage map[string]userModel.User
}

func NewStorageUserMemory() (*StorageUserMemory, error) {
	return &StorageUserMemory{
		storage: make(map[string]userModel.User),
	}, nil
}

func (su *StorageUserMemory) IsUserExists(user userModel.UserInput) (int, error) {
	su.mx2.RLock()
	defer su.mx2.RUnlock()

	if val, ok := su.storage[user.Name]; ok {
		if err := bcrypt.CompareHashAndPassword([]byte(val.Password), []byte(user.Password)); err != nil {
			return customErrors.WRONG_PASSWORD_ERROR, nil
		}
		return val.Id, nil
	}

	return customErrors.NO_USER_ERROR, nil
}

func (su *StorageUserMemory) AddUser(newUser userModel.User) (int, error) {
	su.mx1.Lock()
	defer su.mx1.Unlock()

	newUserInput := userModel.UserInput{
		Name: newUser.Name,
		Password: newUser.Password,
	}

	id, err := su.IsUserExists(newUserInput)

	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	if id != customErrors.NO_USER_ERROR {
		return customErrors.USER_EXISTS_ERROR, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	newUser.Password = string(passwordHash)

	newId := len(su.storage) + 1
	su.storage[newUser.Name] = newUser
	return newId, nil
}
