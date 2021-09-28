package impl

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	userModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[int]userModel.User
}

func NewStorageUserMemory() (*StorageUserMemory, error) {
	return &StorageUserMemory{
		storage: make(map[int]userModel.User),
	}, nil
}

func (su *StorageUserMemory) IsUserExists(user userModel.UserInput) (int, error) {
	su.mx.RLock()
	defer su.mx.RUnlock()

	for key, val := range su.storage {
		if val.Name == user.Name {
			if err := bcrypt.CompareHashAndPassword([]byte(val.Password), []byte(user.Password)); err != nil {
				return customErrors.WRONG_PASSWORD_ERROR, nil
			}
			return key, nil
		}
	}
	return customErrors.NO_USER_ERROR, nil
}

func (su *StorageUserMemory) AddUser(newUser userModel.User) (int, error) {
	su.mx.Lock()
	defer su.mx.Unlock()

	for _, val := range su.storage {
		if val == newUser {
			return customErrors.USER_EXISTS_ERROR, nil
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	newUser.Password = string(passwordHash)

	newId := len(su.storage) + 1
	su.storage[newId] = newUser
	return newId, nil
}
