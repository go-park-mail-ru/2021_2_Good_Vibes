package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type StorageBasketPostgres struct {
	db *sql.DB
}

func NewStorageBasketDB(db *sql.DB, err error) (*StorageBasketPostgres, error) {
	if err != nil {
		return nil, err
	}

	return &StorageBasketPostgres{
		db: db,
	}, nil
}

func (sb *StorageBasketPostgres) PutInBasket(basketProduct models.BasketProduct) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"insert into basket (user_id) values ($1) on conflict(user_id) do nothing",
			basketProduct.UserId,
		)

		if err != nil {
			return err
		}

		_, err = tx.Exec(
			"insert into basket_products (user_id, product_id, count) values ($1, $2, $3)",
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

func (sb *StorageBasketPostgres) DropBasket(userId int) error {
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

func (sb *StorageBasketPostgres) DeleteProduct(product models.BasketProduct) error {
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
			if rows.Err() != nil {
				return rows.Err()
			}
			_, err := tx.Exec(`delete from basket where user_id=$1`, product.UserId)
			if err != nil {
				return err
			}
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
