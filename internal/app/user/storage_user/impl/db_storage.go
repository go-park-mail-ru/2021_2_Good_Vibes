package impl

import (
	"context"
	"errors"
	user_model "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/user"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type StorageUserDB struct {
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

func (su *StorageUserDB) IsUserExists(user user_model.UserInput) (int, error) {
	rows, err := su.conn.Query(context.Background(), "SELECT * FROM customers WHERE name=$1", user.Name)

	if err != nil {
		return -1, err
	}

	defer rows.Close()

	if !rows.Next() {
		return -1, nil
	}

	var tmp user_model.User
	var id int

	err = rows.Scan(&id, &tmp.Name, &tmp.Email, &tmp.Password)
	if err != nil {
		return -1, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(user.Password)); err != nil {
		err = errors.New("wrong password")
		return -1, err
	}

	if rows.Err() != nil {
		return -1, err
	}

	return id, nil
}

func (su *StorageUserDB) AddUser(newUser user_model.User) (int, error) {
	user := user_model.UserInput{
		Name:     newUser.Name,
		Password: newUser.Password,
	}

	id, err := su.IsUserExists(user)

	if err != nil {
		return -1, err
	}

	if id != -1 {
		err := errors.New("user exists")
		return -1, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return -1, err
	}

	rows := su.conn.QueryRow(context.Background(), "INSERT INTO customers (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		newUser.Name,
		newUser.Email,
		passwordHash)

	err = rows.Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
