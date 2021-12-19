package posgresql

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type StorageBrandPostgres struct {
	db *sql.DB
}

func NewStorageBrandDB(db *sql.DB, err error) (*StorageBrandPostgres, error) {
	if err != nil {
		return nil, err
	}
	return &StorageBrandPostgres{
		db: db,
	}, nil
}

func (sc *StorageBrandPostgres) GetBrands() ([]models.Brand, error) {
	var brands []models.Brand
	rows, err := sc.db.Query(`select id, name, image from brands order by id`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		brand := models.Brand{}

		err := rows.Scan(&brand.Id, &brand.Name, &brand.Image)
		if err != nil {
			return nil, err
		}
		fmt.Println(brand.Image)
		brands = append(brands, brand)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return brands, nil
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
