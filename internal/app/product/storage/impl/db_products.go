package impl

import (
	"database/sql"
	productModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product"
	"sync"
)

type StorageProductsDB struct {
	mx sync.RWMutex
	db *sql.DB
}

func NewStorageProductsDB(db *sql.DB, err error) (*StorageProductsDB, error) {
	if err != nil {
		return nil, err
	}

	return &StorageProductsDB{
		db: db,
	}, nil
}

func (ph *StorageProductsDB) GetAllProducts() ([]productModel.Product, error)  {
	rows, err := ph.db.Query("select id, name, category_id from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []productModel.Product
	for rows.Next() {
		product := productModel.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return nil, nil
}

func (ph *StorageProductsDB) GetProductById(id int) (productModel.Product, error)  {
	product := productModel.Product{}

	row := ph.db.QueryRow("select id, name, category_id from products where id=$1", id)

	err := row.Scan(&product.Id, &product.Name, &product.Category)
	if err == sql.ErrNoRows {
		return productModel.Product{}, nil
	}
	if err != nil {
		return product, err
	}

	return product, nil
}

func (ph *StorageProductsDB) GetProductsByCategory(category string) ([]productModel.Product, error) {
	var products []productModel.Product
	rows, err := ph.db.Query("select p.id, p.name, nc1.name from products as p " +
		"join categories as nc1 on p.category = nc1.id " +
		"join categories as nc2 on nc1.lft >= nc2.lft AND " +
		"nc1.rgt <= nc2.rgt where nc2.name = $1 order by nc1.id", category)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		product := productModel.Product{}

		err := rows.Scan(&product.Id, &product.Name, &product.Category)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil
}

func (ph *StorageProductsDB) AddProduct(product productModel.Product) (int, error) {
	var lastInsertId int64

	err := ph.db.QueryRow(
		"with a(id) as (select id from categories where name=$5) " +
			"insert into products (image, name, price, rating, category_id) values ($1, $2, $3, $4, (select id from a)) returning id",
		product.Image,
		product.Name,
		product.Price,
		product.Rating,
		product.Category,
	).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}
