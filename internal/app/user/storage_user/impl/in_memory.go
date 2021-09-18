package impl

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user/storage_user"
	"sync"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[int]storage_user.User
}

func NewStorageUserMemory() *StorageUserMemory {
	return &StorageUserMemory{
		storage: make(map[int]storage_user.User),
	}
}

func (su *StorageUserMemory) IsUserExists(user storage_user.UserInput) (int, error) {
	for key, val := range su.storage {
		if val.Name == user.Name && val.Password == user.Password {
			return key, nil
		}
	}
	return -1, nil
}

func (su *StorageUserMemory) AddUser(newUser storage_user.User) (int, error) {
	su.mx.RLock()
	for _, val := range su.storage {
		if val == newUser {
			return -1, nil
		}
	}
	newId := len(su.storage) + 1
	su.storage[newId] = newUser
	su.mx.RUnlock()
	return newId, nil
}
