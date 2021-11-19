package postgresql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageProductsDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "image", "name", "price", "rating", "category_id", "count_in_stock", "description"})

	expect := []models.Product{
		{Id: 1, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "1", CountInStock: 1, Description: "product"},
		{Id: 2, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "2", CountInStock: 1, Description: "product"},
		{Id: 3, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "3", CountInStock: 1, Description: "product"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Image, item.Name, item.Price, item.Rating, item.Category, item.CountInStock, item.Description)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs().
		WillReturnRows(rows)

	result, err := storage.GetAll()
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
		WithArgs().
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetAll()
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
		NewRows([]string{"id", "image", "name"})

	expect = []models.Product{
		{Id: 1, Image: "product", Name: "Phone"},
		{Id: 2, Image: "product", Name: "Phone"},
		{Id: 3, Image: "product", Name: "Phone"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Image, item.Name)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs().
		WillReturnRows(rows)

	_, err = storage.GetAll()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageProductsDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "image", "name", "price", "rating", "category_id", "count_in_stock", "description"})

	expect := models.Product{
		Id: 1, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "1", CountInStock: 1, Description: "product",
	}

	rows.AddRow(expect.Id, expect.Image, expect.Name, expect.Price, expect.Rating, expect.Category, expect.CountInStock, expect.Description)

	mock.
		ExpectQuery("select ...").
		WithArgs(expect.Id).
		WillReturnRows(rows)

	result, err := storage.GetById(expect.Id)
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
		WithArgs(expect.Id).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetById(expect.Id)
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
		NewRows([]string{"id", "image", "name"})

	expect = models.Product{
		Id: 1, Image: "product", Name: "Phone",
	}

	rows.AddRow(expect.Id, expect.Image, expect.Name)

	mock.
		ExpectQuery("select ...").
		WithArgs().
		WillReturnRows(rows)

	_, err = storage.GetAll()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestGetByCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageProductsDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	category := "ALL"

	// good query
	rows := sqlmock.
		NewRows([]string{"id", "image", "name", "price", "rating", "category_id", "count_in_stock", "description"})

	expect := []models.Product{
		{Id: 1, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "1", CountInStock: 1, Description: "product"},
		{Id: 2, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "2", CountInStock: 1, Description: "product"},
		{Id: 3, Image: "product", Name: "Phone", Price: 1000.00, Rating: 5.0, Category: "3", CountInStock: 1, Description: "product"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Image, item.Name, item.Price, item.Rating, item.Category, item.CountInStock, item.Description)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs(category).
		WillReturnRows(rows)

	result, err := storage.GetByCategory(category)
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
		WithArgs(category).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.GetByCategory(category)
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
		NewRows([]string{"id", "image", "name"})

	expect = []models.Product{
		{Id: 1, Image: "product", Name: "Phone"},
		{Id: 2, Image: "product", Name: "Phone"},
		{Id: 3, Image: "product", Name: "Phone"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Image, item.Name)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs(category).
		WillReturnRows(rows)

	_, err = storage.GetByCategory(category)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageProductsDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	product := models.Product{
		Image:        "product",
		Name:         "product",
		Price:        1000.0,
		Rating:       5.0,
		Category:     "1",
		CountInStock: 1000,
		Description:  "product",
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"id"}).AddRow(1)

	mock.
		ExpectQuery("insert into products").
		WithArgs(product.Image, product.Name, product.Price, product.Rating, product.Category, product.CountInStock, product.Description).
		WillReturnRows(rows)

	_, err = storage.Insert(product)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// query error
	mock.
		ExpectQuery("insert into products").
		WithArgs(product.Image, product.Name, product.Price, product.Rating, product.Category, product.CountInStock, product.Description).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = storage.Insert(product)
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
		NewRows([]string{"id", "name"}).AddRow(1, "product")

	mock.
		ExpectQuery("insert into products").
		WithArgs(product.Image, product.Name, product.Price, product.Rating, product.Category, product.CountInStock, product.Description).
		WillReturnRows(rows)

	_, err = storage.Insert(product)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestSaveProductImageName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewStorageProductsDB(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	productId := 1
	fileName := "avatar"

	//ok query

	mock.
		ExpectExec("update products").
		WithArgs(productId, fileName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveProductImageName(productId, fileName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// query error 1

	mock.
		ExpectExec("update products").
		WithArgs(productId, fileName).
		WillReturnError(errors.Errorf("db error"))

	err = storage.SaveProductImageName(productId, fileName)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestNewStorageProductDB_Fail(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	err = errors.New("This is error: ")

	_, err = NewStorageProductsDB(db, err)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
