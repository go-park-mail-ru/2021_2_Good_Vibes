package postgresql

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"testing"
)

func TestPutOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageOrderDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	orderProducts := []models.OrderProducts{
		{1, 5, 1},
	}

	order := models.Order{
		OrderId:  1,
		UserId:   1,
		Date:     "2014-04-04 18:32:59",
		Address:  "Moscow",
		Cost:     2000.0,
		Status:   "new",
		Products: orderProducts,
	}
	//ok query

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into order_products").
		WithArgs(order.Products[0].OrderId, order.Products[0].ProductId, order.Products[0].Number).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("delete from basket").WithArgs(order.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("delete from basket_products").WithArgs(order.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	id, err := storage.PutOrder(order)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id != order.OrderId {
		t.Errorf("bad id: want %v, have %v", order.OrderId, id)
		return
	}

	// query error 1

	mock.ExpectBegin()

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnError(fmt.Errorf("db error"))

	mock.ExpectRollback()

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 2

	mock.ExpectBegin()

	rows = sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into order_products").
		WithArgs(order.Products[0].OrderId, order.Products[0].ProductId, order.Products[0].Number).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("db error")))

	mock.ExpectRollback()

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 3

	mock.ExpectBegin()

	rows = sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into order_products").
		WithArgs(order.Products[0].OrderId, order.Products[0].ProductId, order.Products[0].Number).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("delete from basket").WithArgs(order.UserId).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("db error")))

	mock.ExpectRollback()

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 4

	mock.ExpectBegin()

	rows = sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into order_products").
		WithArgs(order.Products[0].OrderId, order.Products[0].ProductId, order.Products[0].Number).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("delete from basket").WithArgs(order.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.
		ExpectExec("delete from basket_products").WithArgs(order.UserId).
		WillReturnError(fmt.Errorf("db error"))

	mock.ExpectRollback()

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// result error
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"id", "address"}).
		AddRow(order.OrderId, order.Address)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Address, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.ExpectRollback()

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// no products

	orderProducts = nil

	order = models.Order{
		OrderId:  1,
		UserId:   1,
		Date:     "2014-04-04 18:32:59",
		Address:  "Moscow",
		Cost:     2000.0,
		Status:   "new",
		Products: orderProducts,
	}

	_, err = storage.PutOrder(order)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestNewStorageOrderDB_Fail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	err = errors.New("This is error: ")

	_, err = NewStorageOrderDB(db, err)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
