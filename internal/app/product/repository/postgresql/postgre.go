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
		err = rows.Scan(&product.Id, &product.Image, &product.Name, &product.Price, &product.Rating, &product.Category, &product.CountInStock, &product.Description)
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

	err := row.Scan(&product.Id, &product.Image, &product.Name, &product.Price, &product.Rating, &product.Category, &product.CountInStock, &product.Description)
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
		"nc1.rgt <= nc2.rgt where nc2.name = $1 orders by nc1.id", category)
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

func (ph *StorageProductsDB) AddFavouriteProduct(product models.FavouriteProduct) error {
	_, err := ph.db.Exec(
		`insert into favourite_prod (user_id, product_id) values ($1, $2)`,
		product.UserId,
		product.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ph *StorageProductsDB) DeleteFavouriteProduct(product models.FavouriteProduct) error {
	_, err := ph.db.Exec(
		`delete from favourite_prod where user_id=$1 and product_id=$2`,
		product.UserId,
		product.Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ph *StorageProductsDB) GetFavouriteProducts(userId int) ([]models.Product, error) {
	var products []models.Product

	rows, err := ph.db.Query("select p.id, p.name, p.price, p.rating, p.category_id, "+
		"p.count_in_stock, p.description, p.image from products as p "+
		"join favourite_prod fp on p.id = fp.product_id "+
		"where fp.user_id=$1 "+
		"orders by name", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.Product{}

		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Rating,
			&product.Category, &product.CountInStock, &product.Description, &product.Image)
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
		"with a(id) as (select id from categories where name=$4) "+
			"insert into products (name, price, rating, category_id, count_in_stock, description) values ($1, $2, $3, (select id from a), $5, $6) returning id",
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

func (ph *StorageProductsDB) SaveProductImageName(productId int, fileName string) error {
	_, err := ph.db.Exec(`UPDATE products SET image = $2 WHERE id = $1`, productId, fileName)
	if err != nil {
		return err
	}
	return nil
}

func tx(db *sql.DB, fb func(tx *sql.Tx) error) error {
	trx, _ := db.Begin()
	err := fb(trx)
	if err != nil {
		trx.Rollback()
		return err
	}
	trx.Commit()
	return nil
}
