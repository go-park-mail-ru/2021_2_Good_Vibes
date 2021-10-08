package memory

import (
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	guuid "github.com/google/uuid"
	"sync"
)

type StorageUserMemory struct {
	mx      sync.RWMutex
	storage map[string]models.UserDataStorage
}

func NewStorageUserMemory() (*StorageUserMemory, error) {
	return &StorageUserMemory{
		storage: make(map[string]models.UserDataStorage),
	}, nil
}

func (su *StorageUserMemory) GetUserDataByName(name string) (*models.UserDataStorage, error) {
	su.mx.RLock()
	defer su.mx.RUnlock()

	if _, ok := su.storage[name]; ok {
		returnData := su.storage[name]
		return &returnData, nil
	}

	return nil, nil
}

func (su *StorageUserMemory) InsertUser(newUser models.UserDataForReg) (int, error) {

	su.mx.Lock()
	newId := int(guuid.New().ID())

	su.storage[newUser.Name] = models.UserDataStorage{
		Id:       newId,
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
	}
	su.mx.Unlock()

	return newId, nil
}
