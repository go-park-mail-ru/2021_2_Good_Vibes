package posgresql

import (
	"database/sql"
)

type StorageEmailPostgres struct {
	db *sql.DB
}

func NewStorageEmailDB(db *sql.DB, err error) (*StorageEmailPostgres, error) {
	if err != nil {
		return nil, err
	}
	return &StorageEmailPostgres{
		db: db,
	}, nil
}

func (sc *StorageEmailPostgres) ConfirmEmail(email string) error {
	_, err := sc.db.Exec(`update email_confirm set status = 1 where email = $1`, email)

	if err != nil {
		return err
	}

	return nil
}

func (sc *StorageEmailPostgres) GetToken(email string) (string, error) {
	row := sc.db.QueryRow(`select token from email_confirm where email = $1`, email)

	var token string
	err := row.Scan(&token)

	if err == sql.ErrNoRows {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return token,  nil
}
