package memory

import (
	"errors"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"sync"
)

const countProductsOnPage = 2

type StorageProductsMemory struct {
	mx      sync.RWMutex
	storage map[int]models.Product
}

func NewStorageProductsMemory() (*StorageProductsMemory, error) {
	return &StorageProductsMemory{
		storage: make(map[int]models.Product),
	}, nil
}

func (sp *StorageProductsMemory) Insert(prod models.Product) error {
	sp.mx.Lock()
	defer sp.mx.Unlock()

	newId := prod.Id
	sp.storage[newId] = prod
	return nil
}

func (sp *StorageProductsMemory) GetAll() ([]models.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	var result []models.Product
	for i := 1; i < len(sp.storage)+1; i++ {
		result = append(result, sp.storage[i])
	}
	return result, nil
}

func (sp *StorageProductsMemory) GetById(id int) (models.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	result, ok := sp.storage[id]
	if !ok {
		err := errors.New("product does not exist")
		return result, err
	}
	return result, nil
}

func (sp *StorageProductsMemory) GetOnPage(page int) ([]models.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	result := make([]models.Product, 0, countProductsOnPage)

	startGettingProductsId := countProductsOnPage*page + 1
	for i := startGettingProductsId; i < startGettingProductsId+countProductsOnPage; i++ {
		result = append(result, sp.storage[i])
	}

	return result, nil
}

//чтоб под интерфейс подходило
func (sp *StorageProductsMemory) GetByCategory(category string) ([]models.Product, error) {
	return nil, nil
}
