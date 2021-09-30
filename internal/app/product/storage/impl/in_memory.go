package impl

import (
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"sync"
)

const countProductsOnPage = 2

type StorageProductsMemory struct {
	mx      sync.RWMutex
	storage map[int]product.Product
}

func NewStorageProductsMemory() *StorageProductsMemory {
	return &StorageProductsMemory{
		storage: make(map[int]product.Product),
	}
}

func (sp *StorageProductsMemory) AddProduct(prod product.Product) error {
	sp.mx.Lock()
	defer sp.mx.Unlock()

	newId := prod.Id
	sp.storage[newId] = prod
	return nil
}

func (sp *StorageProductsMemory) GetAllProducts() ([]product.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	var result []product.Product
	for i := 1; i < len(sp.storage)+1; i++ {
		result = append(result, sp.storage[i])
	}
	return result, nil
}

func (sp *StorageProductsMemory) GetProductById(id int) (product.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	result, ok := sp.storage[id]
	if !ok {
		err := errors.New("product does not exist")
		return result, err
	}
	return result, nil
}

func (sp *StorageProductsMemory) GetProductsOnPage(page int) ([]product.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	result := make([]product.Product, 0, countProductsOnPage)

	startGettingProductsId := countProductsOnPage*page + 1
	for i := startGettingProductsId; i < startGettingProductsId+countProductsOnPage; i++ {
		result = append(result, sp.storage[i])
	}

	return result, nil
}
