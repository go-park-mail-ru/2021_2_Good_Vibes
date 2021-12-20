package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"strings"
	"time"
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
	rows, err := ph.db.Query(`select id, image, name, price, rating, category_id, count_in_stock, description, sales, sales_price from products order by id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name,
			            &product.Price, &product.Rating, &product.Category,
			            &product.CountInStock, &product.Description,
			            &product.Sales, &product.SalesPrice)
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


func (ph *StorageProductsDB) GetNewProducts() ([]models.Product, error) {
	layoutISO := "2006-01-02"
	date := time.Now().Format(layoutISO)
	fmt.Println(date)
	rows, err := ph.db.Query(`select id, image, name, price, rating, category_id, count_in_stock, description, sales, sales_price from products where $1::date < '3 day'::interval + date_created order by id desc`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name,
			&product.Price, &product.Rating, &product.Category,
			&product.CountInStock, &product.Description,
			&product.Sales, &product.SalesPrice)
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

func (ph *StorageProductsDB) GetSalesProducts() ([]models.Product, error) {
	rows, err := ph.db.Query("select id, image, name, price, rating, category_id, " +
		                           "count_in_stock, description, sales, sales_price from products " +
		                           "where sales = true order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name,
			&product.Price, &product.Rating, &product.Category,
			&product.CountInStock, &product.Description,
			&product.Sales, &product.SalesPrice)
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

func (ph *StorageProductsDB) GetProductsByBrand(brandId int) ([]models.Product, error) {
	rows, err := ph.db.Query("select p.id, p.image, p.name, p.price, p.rating, p.category_id, " +
		                           "p.count_in_stock, p.description, p.sales, p.sales_price from brands as b " +
	                               "join products p on b.id = p.brand_id " +
		                           "where b.id=$1 order by p.id", brandId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name,
			&product.Price, &product.Rating, &product.Category,
			&product.CountInStock, &product.Description,
			&product.Sales, &product.SalesPrice)
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

	row := ph.db.QueryRow("select id, image, name, price, rating, category_id, count_in_stock, description, sales, sales_price from products where id=$1", id)

	err := row.Scan(&product.Id, &product.Image, &product.Name,
		            &product.Price, &product.Rating, &product.Category,
		            &product.CountInStock, &product.Description,
	             	&product.Sales, &product.SalesPrice)
	if err == sql.ErrNoRows {
		return models.Product{}, nil
	}
	if err != nil {
		return product, err
	}

	return product, nil
}

func (ph *StorageProductsDB) GetByCategory(filter postgre.Filter) ([]models.Product, error) {
	var products []models.Product
	if filter.OrderBy != postgre.TypeOrderRating && filter.OrderBy != postgre.TypeOrderPrice {
		filter.OrderBy = postgre.TypeOrderRating
	}
	if filter.TypeOrder != postgre.TypeOrderMin && filter.TypeOrder != postgre.TypeOrderMax {
		filter.TypeOrder = postgre.TypeOrderMin
	}
	filter.OrderBy = "p." + filter.OrderBy

	rows, err := ph.db.Query("select p.id, p.image, p.name, p.price, p.rating, nc1.name, p.count_in_stock, p.description, p.sales, p.sales_price from products as p "+
		"join categories as nc1 on p.category_id = nc1.id "+
		"join categories as nc2 on nc1.lft >= nc2.lft AND "+
		"nc1.rgt <= nc2.rgt "+
		"where nc2.name = $1 and p.price >= $2 and p.price <= $3 "+
		"and p.rating >= $4 and p.rating <= $5 "+
		"order by "+filter.OrderBy+" "+filter.TypeOrder, filter.NameCategory, filter.MinPrice,
		filter.MaxPrice, filter.MinRating, filter.MaxRating)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		product := models.Product{}

		err := rows.Scan(&product.Id, &product.Image, &product.Name,
			             &product.Price, &product.Rating, &product.Category,
			             &product.CountInStock, &product.Description,
			             &product.Sales, &product.SalesPrice)
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

func (ph *StorageProductsDB) IsFavourite(productID int, userID int64) (*bool, error) {
	var productIDFromDB int
	boolPointer := new(bool)
	*boolPointer = false
	err := ph.db.QueryRow(`select product_id from favourite_prod where product_id=$1 and user_id=$2`, productID, userID).
		Scan(&productIDFromDB)
	if err == sql.ErrNoRows {
		return boolPointer, nil
	}
	if err != nil {
		return boolPointer, err
	}

	*boolPointer = true

	return boolPointer,  nil
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
		"p.count_in_stock, p.description, p.image, p.sales, p.sales_price from products as p "+
		"join favourite_prod fp on p.id = fp.product_id "+
		"where fp.user_id=$1 "+
		"order by name", userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.Product{}

		err := rows.Scan(&product.Id, &product.Name, &product.Price,
						 &product.Rating, &product.Category, &product.CountInStock,
						 &product.Description, &product.Image,
						 &product.Sales, &product.SalesPrice)
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
	product.SalesPrice = product.Price
	var lastInsertId int64

	err := ph.db.QueryRow(
		"with a(id) as (select id from categories where name=$4), "+
			"b(id) as (select id from brands where name=$7) " +
			"insert into products (name, price, rating, category_id, " +
			"count_in_stock, description, brand_id, date_created, sales_price) " +
			"values ($1, $2, $3, (select id from a), $5, $6, (select id from b), now(), $8) returning id",
		&product.Name,
		&product.Price,
		&product.Rating,
		&product.Category,
		&product.CountInStock,
		&product.Description,
		&product.BrandName,
		&product.SalesPrice,
		//TODO: добавить новые поля в инзерт продукта
	).Scan(&lastInsertId)

	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func (ph *StorageProductsDB) SaveProductImageName(productId int, fileName string) error {
	_, err := ph.db.Exec(`UPDATE products SET image = image || ';' || $2 WHERE id = $1`, productId, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (ph *StorageProductsDB) PutSalesProduct(sales models.SalesProduct) error {
	_, err := ph.db.Exec(`update products set sales=true, sales_price=$2 where id = $1`,
		sales.ProductId, sales.SalesPrice)
	if err != nil {
		return err
	}
	return nil
}

func (ph *StorageProductsDB) ChangeRecommendUser(userId int, ProductId int, isSearch string) error {
	var id, userIdGet, productId, counter, counterCurrent int

	if isSearch == "true" {
		counter = 3
	} else {
		counter = 1
	}

	row := ph.db.QueryRow("SELECT id, user_id, product_id, counter FROM recommendation"+
		" WHERE user_id = $1 and product_id = $2", userId, ProductId)
	err := row.Scan(&id, &userIdGet, &productId, &counterCurrent)

	if err == sql.ErrNoRows {
		err = ph.db.QueryRow("INSERT INTO recommendation (user_id, product_id, counter) "+
			"values ($1, $2, $3) returning id", userId, ProductId, counter).Scan(&id)
		if err != nil {
			return err
		}
	} else {
		counterCurrent += counter
		_, err := ph.db.Exec(`UPDATE recommendation SET user_id = $1, product_id = $2,
                          counter = $3 WHERE id = $4`, userIdGet, productId, counterCurrent, id)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (ph *StorageProductsDB)TryGetProductWithSimilarName(productName string) ([]models.Product, error) {
	var searchStr strings.Builder
	searchStr.WriteString(productName)
	searchStr.WriteRune('%')
	rows, err := ph.db.Query("select id, name from products where name ilike $1 " +
		"order by rating desc limit 1", searchStr)
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
