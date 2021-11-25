package postgresql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"testing"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageUserDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	newUser := models.UserDataForReg{
		Name:     "BUSH",
		Email:    "BUSH@mail.ru",
		Password: "12345",
	}

	// good query

	rows := sqlmock.
		NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery("INSERT INTO customers").
		WithArgs(newUser.Name, newUser.Email, newUser.Password).
		WillReturnRows(rows)

	result, err := storage.InsertUser(newUser)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if result != 1 {
		t.Errorf("results not match, want %v, have %v", 1, result)
		return
	}

	// query error
	mock.
		ExpectQuery("INSERT INTO customers").
		WithArgs(newUser.Name, newUser.Email, newUser.Password).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.InsertUser(newUser)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.
		NewRows([]string{"id", "number"}).AddRow(1, 1)

	mock.
		ExpectQuery("INSERT INTO customers").
		WithArgs(newUser.Name, newUser.Email, newUser.Password).
		WillReturnRows(rows)

	_, err = storage.InsertUser(newUser)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
