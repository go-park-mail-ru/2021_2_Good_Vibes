package postgresql

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"reflect"
	"testing"
)

func TestPutOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewOrderRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	orderProducts := []models.OrderProducts{
		{1, 5, 1},
	}

	address := models.Address{
		Country: "Russia",
		Region:  "Moscow",
		City:    "Moscow",
		Street:  "Izmailovskiy prospect",
		House:   "73B",
		Flat:    "44",
		Index:   "109834",
	}

	order := models.Order{
		OrderId:  1,
		UserId:   1,
		Date:     "2014-04-04 18:32:59",
		Address:  address,
		Cost:     2000.0,
		Status:   "new",
		Products: orderProducts,
	}
	//ok query

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
		WillReturnRows(rows)
	mock.
		ExpectExec("insert into delivery_address").
		WithArgs(
			order.OrderId,
			order.Address.Country,
			order.Address.Region,
			order.Address.City,
			order.Address.Street,
			order.Address.House,
			order.Address.Flat,
			order.Address.Index).WillReturnResult(sqlmock.NewResult(1, 1))

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
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
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
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into delivery_address").
		WithArgs(
			order.OrderId,
			order.Address.Country,
			order.Address.Region,
			order.Address.City,
			order.Address.Street,
			order.Address.House,
			order.Address.Flat,
			order.Address.Index).
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
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into delivery_address").
		WithArgs(
			order.OrderId,
			order.Address.Country,
			order.Address.Region,
			order.Address.City,
			order.Address.Street,
			order.Address.House,
			order.Address.Flat,
			order.Address.Index).WillReturnResult(sqlmock.NewResult(1, 1))

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

	// query error 4

	mock.ExpectBegin()

	rows = sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into delivery_address").
		WithArgs(
			order.OrderId,
			order.Address.Country,
			order.Address.Region,
			order.Address.City,
			order.Address.Street,
			order.Address.House,
			order.Address.Flat,
			order.Address.Index).WillReturnResult(sqlmock.NewResult(1, 1))

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

	// query error 5

	mock.ExpectBegin()

	rows = sqlmock.NewRows([]string{"id"}).AddRow(order.OrderId)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
		WillReturnRows(rows)

	mock.
		ExpectExec("insert into delivery_address").
		WithArgs(
			order.OrderId,
			order.Address.Country,
			order.Address.Region,
			order.Address.City,
			order.Address.Street,
			order.Address.House,
			order.Address.Flat,
			order.Address.Index).WillReturnResult(sqlmock.NewResult(1, 1))

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
		AddRow(order.OrderId, order.Cost)

	mock.
		ExpectQuery(`insert into orders`).
		WithArgs(order.UserId, order.Date, order.Cost, order.Status).
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
		Address:  address,
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

func TestSelectPrices(t *testing.T) {
	orderProducts := []models.OrderProducts{
		{1, 5, 1},
	}

	expectedProductPrices := []models.ProductPrice{
		{
			Id:    1,
			Price: 10000,
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewOrderRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	//ok query
	rows := sqlmock.NewRows([]string{"id", "price"}).AddRow(1, 10000)

	mock.
		ExpectQuery(`select`).
		WithArgs().
		WillReturnRows(rows)

	productPrices, err := storage.SelectPrices(orderProducts)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if !reflect.DeepEqual(productPrices, expectedProductPrices) {
		t.Errorf("bad id: want %v, have %v", expectedProductPrices, productPrices)
		return
	}

	// error 1
	mock.
		ExpectQuery(`select`).
		WithArgs().
		WillReturnError(errors.New("new error"))

	_, err = storage.SelectPrices(orderProducts)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err == nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// error 2
	rows = sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery(`select`).
		WithArgs().
		WillReturnRows(rows)

	_, err = storage.SelectPrices(orderProducts)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err == nil {
		t.Errorf("unexpected err: %s", err)
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

	_, err = NewOrderRepository(db, err)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
