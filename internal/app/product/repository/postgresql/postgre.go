package postgresql

import (
	"database/sql"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
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

func (ph *StorageProductsDB) GetAll() ([]models.Product, error) {
	rows, err := ph.db.Query("select id, name, category_id from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return nil, nil
}

func (ph *StorageProductsDB) GetProductById(id int) (models.Product, error) {
	product := models.Product{}

	row := ph.db.QueryRow("select id, name, category_id from products where id=$1", id)

	err := row.Scan(&product.Id, &product.Name, &product.Category)
	if err == sql.ErrNoRows {
		return models.Product{}, nil
	}
	if err != nil {
		return product, err
	}

	return product, nil
}

func (ph *StorageProductsDB) GetByCategory(category string) ([]models.Product, error) {
	var products []models.Product
	rows, err := ph.db.Query("select p.id, p.name, nc1.name from products as p "+
		"join categories as nc1 on p.category = nc1.id "+
		"join categories as nc2 on nc1.lft >= nc2.lft AND "+
		"nc1.rgt <= nc2.rgt where nc2.name = $1 order by nc1.id", category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.Product{}

		err := rows.Scan(&product.Id, &product.Name, &product.Category)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil
}

func (ph *StorageProductsDB) Insert(product models.Product) error {
	var lastInsertId int64

	err := ph.db.QueryRow(
		"with a(id) as (select id from categories where name=$5) "+
			"insert into products (image, name, price, rating, category_id) values ($1, $2, $3, $4, (select id from a)) returning id",
		product.Image,
		product.Name,
		product.Price,
		product.Rating,
		product.Category,
	).Scan(&lastInsertId)

	if err != nil {
		return err
	}

	return nil
}
