package postgresql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestPutInBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewBasketRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	basketProduct := models.BasketProduct{
		UserId:    1,
		ProductId: 3,
		Number:    5,
	}

	//ok query

	mock.ExpectBegin()

	mock.
		ExpectExec("insert into basket").
		WithArgs(
			basketProduct.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("insert into basket_products").
		WithArgs(basketProduct.UserId, basketProduct.ProductId, basketProduct.Number).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = storage.PutInBasket(basketProduct)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.ExpectBegin()

	mock.
		ExpectExec("insert into basket").
		WithArgs(
			basketProduct.UserId).
		WillReturnResult(sqlmock.NewErrorResult(errors.Errorf("db error")))

	mock.ExpectRollback()

	err = storage.PutInBasket(basketProduct)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 2

	mock.ExpectBegin()

	mock.
		ExpectExec("insert into basket").
		WithArgs(
			basketProduct.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("insert into basket_products (.+) values (.+)").
		WithArgs(basketProduct.UserId, basketProduct.ProductId, basketProduct.Number).
		WillReturnError(errors.Errorf("db error"))

	mock.ExpectRollback()

	err = storage.PutInBasket(basketProduct)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestGetBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewBasketRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"product_id", "count"})

	expect := []models.BasketProduct{
		{ProductId: 1, Number: 1},
		{ProductId: 2, Number: 3},
		{ProductId: 3, Number: 5},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.ProductId, item.Number)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := storage.GetBasket(1)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("results not match, want %v, have %v", expect, result)
		return
	}

	// query error
	mock.
		ExpectQuery("select ...").
		WithArgs(1).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetBasket(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"product_id"}).
		AddRow(3)

	mock.
		ExpectQuery("select ...").
		WithArgs(1).
		WillReturnRows(rows)

	_, err = storage.GetBasket(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDropBasket(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewBasketRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	userId := 1

	//ok query

	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket").
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = storage.DropBasket(userId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket").
		WithArgs(userId).
		WillReturnResult(sqlmock.NewErrorResult(errors.Errorf("db error")))

	mock.ExpectRollback()

	err = storage.DropBasket(userId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 2

	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket").
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(userId).
		WillReturnError(errors.Errorf("db error"))

	mock.ExpectRollback()

	err = storage.DropBasket(userId)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewBasketRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	product := models.BasketProduct{
		UserId:    1,
		ProductId: 1,
	}

	//ok query, valid rows
	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.
		NewRows([]string{"user_id", "product_id", "count"}).AddRow(1, 1, 1)

	mock.
		ExpectQuery("select ...").
		WithArgs(product.ProductId).
		WillReturnRows(rows)

	mock.ExpectCommit()

	err = storage.DeleteProduct(product)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	//ok query, no rows
	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows = sqlmock.
		NewRows([]string{"user_id", "product_id", "count"})

	mock.
		ExpectQuery("select ...").
		WithArgs(product.ProductId).
		WillReturnRows(rows)

	mock.
		ExpectExec("delete from basket").
		WithArgs(product.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = storage.DeleteProduct(product)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewErrorResult(errors.Errorf("db error")))

	mock.ExpectRollback()

	err = storage.DeleteProduct(product)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error 2
	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectQuery("select ...").
		WithArgs(product.ProductId).
		WillReturnError(errors.Errorf("db error"))

	mock.ExpectRollback()

	err = storage.DeleteProduct(product)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// scan error

	mock.ExpectBegin()

	rows = sqlmock.
		NewRows([]string{"user_id", "product_id"})

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectQuery("select ...").
		WithArgs(product.ProductId).
		WillReturnRows(rows)

	mock.ExpectRollback()

	err = storage.DeleteProduct(product)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// error 3

	mock.ExpectBegin()

	mock.
		ExpectExec("delete from basket_products").
		WithArgs(product.UserId, product.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows = sqlmock.
		NewRows([]string{"user_id", "product_id", "count"})

	mock.
		ExpectQuery("select ...").
		WithArgs(product.ProductId).
		WillReturnRows(rows)

	mock.
		ExpectExec("delete from basket").
		WithArgs(product.UserId).
		WillReturnError(errors.Errorf("db error"))

	mock.ExpectRollback()

	err = storage.DeleteProduct(product)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestNewStorageBasketDB_Fail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	err = errors.New("This is error: ")

	_, err = NewBasketRepository(db, err)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
