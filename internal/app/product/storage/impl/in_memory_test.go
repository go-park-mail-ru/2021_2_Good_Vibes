package impl

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"sync"
	"testing"
)

func storageInit() map[int]product.Product {
	storage := make(map[int]product.Product)

	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)
	product2 := product.NewProduct(2, "images/cat2.jpeg", "cat2", 1000, 100)

	storage[product1.Id] = product1
	storage[product2.Id] = product2

	return storage
}

func TestStorageProductsMemory_AddProduct(t *testing.T) {
	storage := storageInit()
	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)

	type args struct {
		prod product.Product
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"success",
			args{product1},
			false },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &StorageProductsMemory{
				mx:      sync.RWMutex{},
				storage: storage,
			}
			if err := sp.AddProduct(tt.args.prod); (err != nil) != tt.wantErr {
				t.Errorf("AddProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestStorageProductsMemory_GetAllProducts(t *testing.T) {
	storage := storageInit()

	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)
	product2 := product.NewProduct(2, "images/cat2.jpeg", "cat2", 1000, 100)

	var products []product.Product
	products = append(products, product1)
	products = append(products, product2)

	tests := []struct {
		name    string

		want []product.Product
		wantErr bool
	}{
		{"success",
			products,
			false },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &StorageProductsMemory{
				mx:      sync.RWMutex{},
				storage: storage,
			}
			got, err := sp.GetAllProducts()
			gotJs, _ := json.Marshal(got)
			wantJs, _ := json.Marshal(tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotJs) != string(wantJs) {
				t.Errorf("GetAllProducts() got = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestStorageProductsMemory_GetProductById(t *testing.T) {
	storage := storageInit()

	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)
	product2 := product.NewProduct(0, "", "", 0, 0)

	tests := []struct {
		name    string
		id int
		want product.Product
		wantErr bool
	}{
		{"success",
			1,
			product1,
			false },
		{"fail",
			4,
			product2,
			true },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &StorageProductsMemory{
				mx:      sync.RWMutex{},
				storage: storage,
			}
			got, err := sp.GetProductById(tt.id)
			gotJs, _ := json.Marshal(got)
			wantJs, _ := json.Marshal(tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetProductById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotJs) != string(wantJs) {
				t.Errorf("GetProductById() got = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestStorageProductsMemory_GetProductsOnPage(t *testing.T) {
	storage := storageInit()

	product1 := product.NewProduct(1, "images/cat1.jpeg", "cat1", 1000, 100)
	product2 := product.NewProduct(2, "images/cat2.jpeg", "cat2", 1000, 100)

	var products []product.Product
	products = append(products, product1)
	products = append(products, product2)

	tests := []struct {
		name    string
		id int
		want []product.Product
		wantErr bool
	}{
		{"success",
			0,
			products,
			false },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := &StorageProductsMemory{
				mx:      sync.RWMutex{},
				storage: storage,
			}
			got, err := sp.GetProductsOnPage(tt.id)
			gotJs, _ := json.Marshal(got)
			wantJs, _ := json.Marshal(tt.want)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotJs) != string(wantJs) {
				t.Errorf("GetAllProducts() got = %v, want %v", got, tt.want)
			}
		})
	}
}