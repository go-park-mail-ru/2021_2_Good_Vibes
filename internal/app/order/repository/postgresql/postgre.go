package postgresql

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"strconv"
	"strings"
)

const (
	FieldsNum = 3
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB, err error) (*OrderRepository, error) {
	if err != nil {
		return nil, err
	}

	return &OrderRepository{
		db: db,
	}, nil
}

func (so *OrderRepository) PutOrder(order models.Order) (int, error) {
	if order.Products == nil {
		err := errors.New("No products, ")
		return 0, err
	}

	err := tx(so.db, func(tx *sql.Tx) error {
		err := tx.QueryRow(
			`insert into orders (user_id, date, cost, status) values ($1, $2, $3, $4) returning id`,
			order.UserId,
			order.Date,
			order.Cost,
			order.Status,
		).Scan(&order.OrderId)

		if err != nil {
			return err
		}

		address := order.Address

		_, err = tx.Exec(
			`insert into delivery_address (order_id, country, region, city, street, house, flat, a_index) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
			order.OrderId,
			address.Country,
			address.Region,
			address.City,
			address.Street,
			address.House,
			address.Flat,
			address.Index)

		if err != nil {
			return err
		}

		for i, _ := range order.Products {
			order.Products[i].OrderId = order.OrderId
		}

		query, values := makeOrderProductsInsertQuery(order)

		_, err = tx.Exec(query, values...)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`delete from basket where user_id=$1`, order.UserId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`delete from basket_products where user_id=$1`, order.UserId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return order.OrderId, nil
}

func (so *OrderRepository) SelectPrices(products []models.OrderProducts) ([]models.ProductPrice, error) {
	query := makeSelectPricesQuery(products)

	rows, err := so.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var productPrices []models.ProductPrice

	for rows.Next() {
		productPrice := models.ProductPrice{}
		err = rows.Scan(&productPrice.Id, &productPrice.Price)
		if err != nil {
			return nil, err
		}

		productPrices = append(productPrices, productPrice)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return productPrices, nil
}

func makeSelectPricesQuery(products []models.OrderProducts) string {
	query := strings.Builder{}
	query.WriteString("select id, price from products where id in ")
	query.WriteString(`(`)

	str := make([]string, len(products))

	for i, product := range products {
		str[i] = strconv.Itoa(product.ProductId)
	}

	query.WriteString(strings.Join(str, `,`))
	query.WriteString(`)`)

	return query.String()
}

func makeOrderProductsInsertQuery(order models.Order) (string, []interface{}) {
	query := strings.Builder{}
	query.WriteString("insert into order_products (order_id, product_id, count) values")

	values := make([]interface{}, FieldsNum*len(order.Products))
	for i, s := range order.Products {
		values[i*FieldsNum] = s.OrderId
		values[i*FieldsNum+1] = s.ProductId
		values[i*FieldsNum+2] = s.Number

		n := i * FieldsNum

		query.WriteString(`(`)
		str := make([]string, FieldsNum)
		for j := 0; j < FieldsNum; j++ {
			str[j] = `$` + strconv.Itoa(n+j+1)
		}

		query.WriteString(strings.Join(str, ","))
		query.WriteString(`),`)
	}

	str := query.String()

	return str[:len(str)-1], values
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
