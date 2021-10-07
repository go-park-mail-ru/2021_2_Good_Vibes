package impl

import (
	"context"
	"errors"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	userModel "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"sync"
)

type StorageUserDB struct {
	mx sync.RWMutex
	conn *pgx.Conn
}

func NewStorageUserDB() (*StorageUserDB, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}
	return &StorageUserDB{
			conn: conn,
		},
		nil
}

func (su *StorageUserDB) IsUserExists(user userModel.UserInput) (int, error) {
	su.mx.Lock()
	defer su.mx.Unlock()
	rows, err := su.conn.Query(context.Background(), "SELECT * FROM customers WHERE name=$1", user.Name)

	if err != nil {
		return customErrors.DB_ERROR, err
	}

	defer rows.Close()

	if !rows.Next() {
		return customErrors.NO_USER_ERROR, nil
	}

	var tmp userModel.User
	var id int

	err = rows.Scan(&id, &tmp.Name, &tmp.Email, &tmp.Password)
	if err != nil {
		return customErrors.DB_ERROR, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(user.Password)); err != nil {
		err = errors.New(customErrors.WRONG_PASSWORD_DESCR)
		return customErrors.WRONG_PASSWORD_ERROR, err
	}

	if rows.Err() != nil {
		return customErrors.DB_ERROR, err
	}

	return id, nil
}

func (su *StorageUserDB) AddUser(newUser userModel.User) (int, error) {
	user := userModel.UserInput{
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

	su.mx.Lock()
	defer su.mx.Unlock()

	rows := su.conn.QueryRow(context.Background(), "INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Name,
		newUser.Email,
		passwordHash)

	err = rows.Scan(&id)

	if err != nil {
		return customErrors.DB_ERROR, err
	}

	return id, nil
}
