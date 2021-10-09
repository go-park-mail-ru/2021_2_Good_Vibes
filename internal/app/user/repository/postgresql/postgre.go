package postgresql

import (
	"database/sql"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"golang.org/x/crypto/bcrypt"
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

func (su *StorageUserDB) IsUserExists(user models.UserDataForInput) (int, error) {
	var tmp models.UserDataStorage
	row := su.db.QueryRow("SELECT * FROM customers WHERE name=$1", user.Name)

	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password)
	if err == sql.ErrNoRows {
		return customErrors.NO_USER_ERROR, nil
	}

	if err != nil {
		return customErrors.DB_ERROR, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(user.Password)); err != nil {
		err = errors.New(customErrors.WRONG_PASSWORD_DESCR)
		return customErrors.WRONG_PASSWORD_ERROR, err
	}

	return tmp.Id, nil
}

func (su *StorageUserDB) AddUser(newUser models.UserDataForReg) (int, error) {
	user := models.UserDataForInput{
		Name:     newUser.Name,
		Password: newUser.Password,
	}

	id, err := su.IsUserExists(user)

	if err != nil {
		return id, err
	}

	if id != customErrors.NO_USER_ERROR {
		return customErrors.USER_EXISTS_ERROR, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return customErrors.SERVER_ERROR, err
	}

	rows := su.db.QueryRow("INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Name,
		newUser.Email,
		passwordHash)

	err = rows.Scan(&id)

	if err != nil {
		return customErrors.DB_ERROR, err
	}

	return id, nil
}
