package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type StorageProductsDB struct {
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
	rows, err := ph.db.Query("select id, image, name, price, rating, category_id, count_in_stock, description from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name,  &product.Price, &product.Rating, &product.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return products, nil
}

func (ph *StorageProductsDB) GetById(id int) (models.Product, error) {
	product := models.Product{}

	row := ph.db.QueryRow("select id, image, name, price, rating, category_id, count_in_stock, description from products where id=$1", id)

	err := row.Scan(&product.Id, &product.Image, &product.Name,  &product.Price, &product.Rating, &product.Category, &product.CountInStock, &product.Description)
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
	rows, err := ph.db.Query("select p.id, p.image, p.name, p.price, p.rating, nc1.name, p.count_in_stock, p.description from products as p "+
		"join categories as nc1 on p.category_id = nc1.id "+
		"join categories as nc2 on nc1.lft >= nc2.lft AND "+
		"nc1.rgt <= nc2.rgt where nc2.name = $1 order by nc1.id", category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.Product{}

		err := rows.Scan(&product.Id, &product.Image, &product.Name, &product.Price, &product.Rating, &product.Category, &product.CountInStock, &product.Description)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

func (ph *StorageProductsDB) Insert(product models.Product) (int, error) {
	var lastInsertId int64

	err := ph.db.QueryRow(
		"with a(id) as (select id from categories where name=$5) "+
			"insert into products (image, name, price, rating, category_id, count_in_stock, description) values ($1, $2, $3, $4, (select id from a), $6, $7) returning id",
		product.Image,
		product.Name,
		product.Price,
		product.Rating,
		product.Category,
		product.CountInStock,
		product.Description,
	).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func (ph *StorageProductsDB) SaveProductImageName(productId int,fileName string) error {
	_, err := ph.db.Exec(`UPDATE products SET image = $2 WHERE id = $1`, productId, fileName)
	if err != nil {
		return err
	}
	return nil
}
