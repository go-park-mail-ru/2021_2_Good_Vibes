package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"strconv"
	"strings"
)

const (
	FieldsNum = 4
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
			`insert into orders (user_id, date, cost, status, email) values ($1, $2, $3, $4, $5) returning id`,
			order.UserId,
			order.Date,
			order.Cost,
			order.Status,
			order.Email,
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
	fmt.Println(productPrices)
	return productPrices, nil
}

func (so *OrderRepository) GetAllOrders(user int) ([]models.Order, error) {
	var orders []models.Order

	err := tx(so.db, func(tx *sql.Tx) error {
		rows, err := so.db.Query("select id, user_id, date, cost, status, email from orders where user_id = $1 order by date desc", user)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			order := models.Order{}

			err := rows.Scan(&order.OrderId, &order.UserId, &order.Date, &order.Cost, &order.Status, &order.Email)
			if err != nil {
				return err
			}

			orders = append(orders, order)
		}

		if rows.Err() != nil {
			return nil
		}

		for i, _ := range orders {
			var products []models.OrderProducts
			rows, err := so.db.Query("select o.order_id, o.product_id, o.count, o.price, p.image, p.name, " +
				                           "p.rating, p.description, p.sales from order_products as o " +
				                           "join products p on o.product_id = p.id where order_id = $1",
				                          orders[i].OrderId)
			if err != nil {
				return err
			}

			defer rows.Close()

			for rows.Next() {
				product := models.OrderProducts{}

				err := rows.Scan(&product.OrderId, &product.ProductId, &product.Number, &product.Price,
				                 &product.Image, &product.Name, &product.Rating, &product.Description,
				                 &product.Sales)
				if err != nil {
					return err
				}

				products = append(products, product)
			}
			orders[i].Products = products

			if rows.Err() != nil {
				return nil
			}

			rows.Close()
		}

		for i, _ := range orders {
			var address models.Address
			err := so.db.QueryRow("select country, region, city, street, house, flat, a_index from delivery_address where order_id = $1", orders[i].OrderId).
				Scan(
					&address.Country,
					&address.Region,
					&address.City,
					&address.Street,
					&address.House,
					&address.Flat,
					&address.Index,
				)

			if err != nil {
				return err
			}

			orders[i].Address = address
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (so *OrderRepository) CheckPromoCode(promoCode string) (*models.PromoCode, error) {
	var promoReturn models.PromoCode

	row := so.db.QueryRow("select type, code, value, category_id, product_id, uses_left"+
		" from promocode where code = $1", promoCode)
	var categoryId, productId sql.NullInt32
	err := row.Scan(&promoReturn.Type, &promoReturn.Code, &promoReturn.Value,
		&categoryId, &productId, &promoReturn.UsesLeft)

	promoReturn.ProductId = -1
	promoReturn.CategoryId = -1
	if productId.Valid {
		promoReturn.ProductId = int(productId.Int32)
	}
	if categoryId.Valid {
		promoReturn.CategoryId = int(categoryId.Int32)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &promoReturn, nil
}

func (so *OrderRepository) GetOrderById(orderId int) (models.Order, error) {
	var order models.Order

	err := tx(so.db, func(tx *sql.Tx) error {
		err := so.db.QueryRow("select id, user_id, date, cost, status, email from orders where id = $1", orderId).
			Scan(&order.OrderId, &order.UserId, &order.Date, &order.Cost, &order.Status, &order.Status)

		if err == sql.ErrNoRows {
			return nil
		}

		if err != nil {
			return err
		}

		var products []models.OrderProducts
		rows, err := so.db.Query("select order_id, product_id, count from order_products where order_id = $1", orderId)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			product := models.OrderProducts{}

			err := rows.Scan(&product.OrderId, &product.ProductId, &product.Number)
			if err != nil {
				return err
			}

			products = append(products, product)
		}
		order.Products = products

		if rows.Err() != nil {
			return nil
		}

		return nil
	})
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
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
	query.WriteString("insert into order_products (order_id, product_id, count, price) values")

	values := make([]interface{}, FieldsNum*len(order.Products))
	for i, s := range order.Products {
		values[i*FieldsNum] = s.OrderId
		values[i*FieldsNum+1] = s.ProductId
		values[i*FieldsNum+2] = s.Number
		values[i*FieldsNum+3] = s.Price
		fmt.Println(s.Price)
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

func (so *OrderRepository) GetProductCategory(productId int) (int, error) {
	row := so.db.QueryRow("select category_id "+
		" from products where id = $1", productId)
	var categoryId int
	err := row.Scan(&categoryId)
	return categoryId, err
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
