package impl

import (
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	userModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	guuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[string]userModel.UserStorage
}

func NewStorageUserMemory() (*StorageUserMemory, error) {
	return &StorageUserMemory{
		storage: make(map[string]userModel.UserStorage),
	}, nil
}

func (su *StorageUserMemory) IsUserExists(user userModel.UserInput) (int, error) {
	su.mx.RLock()
	defer su.mx.RUnlock()

	if val, ok := su.storage[user.Name]; ok {
		if err := bcrypt.CompareHashAndPassword([]byte(val.Password), []byte(user.Password)); err != nil {
			return customErrors.WRONG_PASSWORD_ERROR, nil
		}
		return val.Id, nil
	}

	return customErrors.NO_USER_ERROR, nil
}

func (su *StorageUserMemory) AddUser(newUser userModel.User) (int, error) {
	id, err := su.IsUserExists(userModel.UserInput{Name: newUser.Name, Password: newUser.Password})

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

	su.mx.Lock()
    newId := int(guuid.New().ID())

	su.storage[newUser.Name] = userModel.UserStorage{
		Id:       newId,
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
	}
	su.mx.Unlock()

	return newId, nil
}
