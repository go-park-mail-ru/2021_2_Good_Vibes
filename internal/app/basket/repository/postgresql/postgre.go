package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type BasketRepository struct {
	db *sql.DB
}

func NewBasketRepository(db *sql.DB, err error) (*BasketRepository, error) {
	if err != nil {
		return nil, err
	}

	return &BasketRepository{
		db: db,
	}, nil
}

func (sb *BasketRepository) PutInBasket(basketProduct models.BasketProduct) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"insert into basket (user_id) values ($1) on conflict(user_id) do nothing",
			basketProduct.UserId,
		)

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			"insert into basket_products (user_id, product_id, count) values ($1, $2, $3) on conflict(user_id,product_id) do update set count=$3",
			basketProduct.UserId,
			basketProduct.ProductId,
			basketProduct.Number,
		)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (sb *BasketRepository) GetBasket(userId int) ([]models.BasketProduct, error) {
	var basketProducts []models.BasketProduct
	rows, err := sb.db.Query("select product_id, count from basket_products where user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.BasketProduct{}

		err := rows.Scan(&product.ProductId, &product.Number)
		if err != nil {
			return nil, err
		}

		basketProducts = append(basketProducts, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return basketProducts, nil
}

func (sb *BasketRepository) DropBasket(userId int) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`delete from basket where user_id=$1`, userId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`delete from basket_products where user_id=$1`, userId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (sb *BasketRepository) DeleteProduct(product models.BasketProduct) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`delete from basket_products where user_id=$1 and product_id=$2`, product.UserId, product.ProductId)
		if err != nil {
			return err
		}

		rows, err := tx.Query(`select from basket_products where user_id=$1`, product.UserId)
		if err != nil {
			return err
		}

		defer rows.Close()
		if !rows.Next() {
			_, err := tx.Exec(`delete from basket where user_id=$1`, product.UserId)
			if err != nil {
				return err
			}
		}

		if rows.Err() != nil {
			return rows.Err()
		}

		return nil
	})

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
