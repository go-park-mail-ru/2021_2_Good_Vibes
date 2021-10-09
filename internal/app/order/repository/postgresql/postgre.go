package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"strconv"
)

type StorageOrderPostgres struct {
	db *sql.DB
}

func NewStorageOrderDB(db *sql.DB, err error) (*StorageOrderPostgres, error) {
	if err != nil {
		return nil, err
	}

	return &StorageOrderPostgres{
		db: db,
	}, nil
}

func (so *StorageOrderPostgres) PutOrder(order models.Order) (int, error) {
	err := so.db.QueryRow(
		"insert into orders (user_id, date, address, cost, status) values ($1, $2, $3, $4, $5) returning id",
		order.UserId,
		order.Date,
		order.Address,
		order.Cost,
		order.Status,
	).Scan(&order.OrderId)

	if err != nil {
		return 0, err
	}

	for i, _ := range order.Products {
		order.Products[i].OrderId = order.OrderId
	}

	query := "insert into order_products (order_id, product_id, count) values"

	var values []interface{}
	for i, s := range order.Products {
		values = append(values, s.OrderId, s.ProductId, s.Number)

		numFields := 3
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$`+strconv.Itoa(n + j + 1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]

	_, err = so.db.Exec(query, values...)

	return order.OrderId, nil
}
