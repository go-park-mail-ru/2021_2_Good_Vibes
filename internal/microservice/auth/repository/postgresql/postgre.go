package postgresql

import (
	"database/sql"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type StorageUserDB struct {
	db *sql.DB
}

func NewStorageUserDB(db *sql.DB, err error) (*StorageUserDB, error) {
	if err != nil {
		return nil, err
	}
	return &StorageUserDB{
		db: db,
	}, nil
}

func (su *StorageUserDB) GetUserDataByName(name string) (*models.UserDataStorage, error) {
	var tmp models.UserDataStorage
	row := su.db.QueryRow("SELECT id, name, email, password FROM customers WHERE name=$1", name)

	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func (su *StorageUserDB) InsertUser(newUser models.UserDataForReg) (int, error) {
	rows := su.db.QueryRow("INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Name,
		newUser.Email,
		newUser.Password)

	var id int
	err := rows.Scan(&id)
	if err != nil {
		return customErrors.DB_ERROR, err
	}

	return id, nil
}
