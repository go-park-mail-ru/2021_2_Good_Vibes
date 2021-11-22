package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestGetUserDataByName(t *testing.T) {
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

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "email", "password"})

	expect := models.UserDataStorage{
		Id: 1, Name: "BUSH", Email: "BUSH@mail.ru", Password: "12345",
	}

	rows.AddRow(expect.Id, expect.Name, expect.Email, expect.Password)

	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Name).
		WillReturnRows(rows)

	result, err := storage.GetUserDataByName(expect.Name)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(*result, expect) {
		t.Errorf("results not match, want %v, have %v", expect, result)
		return
	}

	// query error
	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Name).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetUserDataByName(expect.Name)
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
		NewRows([]string{"id", "name", "email"})

	rows.AddRow(expect.Id, expect.Name, expect.Email)

	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Name).
		WillReturnRows(rows)

	_, err = storage.GetUserDataByName(expect.Name)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

//func TestInsertUser(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	storage, err := NewStorageUserDB(db, nil)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//
//	newUser := models.UserDataForReg{
//		Name:     "BUSH",
//		Email:    "BUSH@mail.ru",
//		Password: "12345",
//	}
//
//	// good query
//
//	rows := sqlmock.
//		NewRows([]string{"id"}).AddRow(1)
//
//	mock.
//		ExpectQuery("insert into customers").
//		WithArgs(newUser.Name, newUser.Email, newUser.Password).
//		WillReturnRows(rows)
//
//	result, err := storage.InsertUser(newUser)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//
//	if result != 1 {
//		t.Errorf("results not match, want %v, have %v", 1, result)
//		return
//	}
//
//	// query error
//	mock.
//		ExpectQuery("insert into customers").
//		WithArgs(newUser.Name, newUser.Email, newUser.Password).
//		WillReturnError(fmt.Errorf("db_error"))
//
//	_, err = storage.InsertUser(newUser)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if err == nil {
//		t.Errorf("expected error, got nil")
//		return
//	}
//
//	// row scan error
//	rows = sqlmock.
//		NewRows([]string{"id", "number"}).AddRow(1, 1)
//
//	mock.
//		ExpectQuery("insert into customers").
//		WithArgs(newUser.Name, newUser.Email, newUser.Password).
//		WillReturnRows(rows)
//
//	_, err = storage.InsertUser(newUser)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if err == nil {
//		t.Errorf("expected error, got nil")
//		return
//	}
//}

func TestGetUserDataById(t *testing.T) {
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

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "name", "email", "password", "avatar"})

	expect := models.UserDataStorage{
		Id: 1, Name: "BUSH", Email: "BUSH@mail.ru", Password: "12345", Avatar: sql.NullString{String: "avatar", Valid: true},
	}

	rows.AddRow(expect.Id, expect.Name, expect.Email, expect.Password, expect.Avatar)

	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Id).
		WillReturnRows(rows)

	result, err := storage.GetUserDataById(uint64(expect.Id))
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(*result, expect) {
		t.Errorf("results not match, want %v, have %v", expect, result)
		return
	}

	// query error
	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Id).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetUserDataById(uint64(expect.Id))
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
		NewRows([]string{"id", "name", "email"})

	rows.AddRow(expect.Id, expect.Name, expect.Email)

	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Id).
		WillReturnRows(rows)

	_, err = storage.GetUserDataById(uint64(expect.Id))
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestSaveAvatarName(t *testing.T) {
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

	userId := 1
	fileName := "avatar"

	//ok query

	mock.
		ExpectExec("update customers").
		WithArgs(userId, fileName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveAvatarName(userId, fileName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.
		ExpectExec("update customers").
		WithArgs(userId, fileName).
		WillReturnError(errors.Errorf("db error"))

	err = storage.SaveAvatarName(userId, fileName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestUpdateUser(t *testing.T) {
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

	user := models.UserDataProfile{
		Id: 1, Name: "BUSH", Email: "BUSH@mail.ru", Avatar: "avatar",
	}

	//ok query

	mock.
		ExpectExec("update customers").
		WithArgs(user.Name, user.Email, user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.UpdateUser(user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.
		ExpectExec("update customers").
		WithArgs(user.Name, user.Email, user.Id).
		WillReturnError(errors.Errorf("db error"))

	err = storage.UpdateUser(user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestNewStorageUserDB_Fail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	err = errors.New("This is error: ")

	_, err = NewStorageUserDB(db, err)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
