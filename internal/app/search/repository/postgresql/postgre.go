package postgresql

import (
	"database/sql"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB, err error) (*SearchRepository, error) {
	if err != nil {
		return nil, err
	}

	return &SearchRepository{
		db: db,
	}, nil
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
