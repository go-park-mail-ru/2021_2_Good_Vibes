package impl

import (
	user_model "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"sync"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[int]user_model.User
}

func NewStorageUserMemory() *StorageUserMemory {
	return &StorageUserMemory{
		storage: make(map[int]user_model.User),
	}
}

func (su *StorageUserMemory) IsUserExists(user user_model.UserInput) (int, error) {
	su.mx.RLock()
	defer su.mx.RUnlock()

	for key, val := range su.storage {
		if val.Name == user.Name && val.Password == user.Password {
			return key, nil
		}
	}
	return -1, nil
}

func (su *StorageUserMemory) AddUser(newUser user_model.User) (int, error) {
	su.mx.Lock()
	defer su.mx.Unlock()

	for _, val := range su.storage {
		if val == newUser {
			return -1, nil
		}
	}
	newId := len(su.storage) + 1
	su.storage[newId] = newUser
	return newId, nil
}
