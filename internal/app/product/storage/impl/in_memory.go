package impl

import (
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"sync"
)

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

	newId := len(sp.storage) + 1
	sp.storage[newId] = prod
	return nil
}

func (sp *StorageProductsMemory) GetAllProducts() ([]product.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()

	var result []product.Product
	for _, value := range sp.storage {
		result = append(result, value)
	}
	return result, nil
}

func (sp *StorageProductsMemory) GetProductsOnPage(page int) ([]product.Product, error) {
	sp.mx.RLock()
	defer sp.mx.RUnlock()
	//это в конфиг, или вообще решить откуда оно браться будет
	countProductsOnPage := 10

	var result []product.Product
	startGettingProductsId := countProductsOnPage*page + 1
	for i := startGettingProductsId; i < startGettingProductsId+countProductsOnPage; i++ {
		result = append(result, sp.storage[i])
	}

	return result, nil
}
