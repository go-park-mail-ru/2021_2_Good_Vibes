package postgresql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"reflect"
	"testing"
)

func Test_GetAllRatingsOfProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewReviewRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	// good query
	rows := sqlmock.
		NewRows([]string{"rating", "count"})

	expect := []models.ProductRating{
		{Rating: 5, Count: 5},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Rating, item.Count)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := storage.GetAllRatingsOfProduct(1)
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

	_, err = storage.GetAllRatingsOfProduct(1)
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
		NewRows([]string{"rating"})

	expect = []models.ProductRating{
		{Rating: 5},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Rating)
	}

	mock.
		ExpectQuery("select ...").
		WithArgs(1).
		WillReturnRows(rows)

	_, err = storage.GetAllRatingsOfProduct(1)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func Test_AddReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewReviewRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	review := models.Review{
		UserId:    1,
		ProductId: 1,
		Rating:    5,
		Text:      "Good",
	}

	var productRating float64

	// good query

	mock.ExpectBegin()
	mock.
		ExpectExec("insert into reviews").
		WithArgs(review.UserId, review.ProductId, review.Rating, review.Text).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("update product").
		WithArgs(productRating, review.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = storage.AddReview(review, productRating)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// query error
	mock.ExpectBegin()
	mock.
		ExpectExec("insert into reviews").
		WithArgs(review.UserId, review.ProductId, review.Rating, review.Text).
		WillReturnError(fmt.Errorf("db_error"))
	mock.ExpectRollback()
	err = storage.AddReview(review, productRating)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func Test_UpdateReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewReviewRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	review := models.Review{
		UserId:    1,
		ProductId: 1,
		Rating:    5,
		Text:      "Good",
	}

	var productRating float64

	// good query

	mock.ExpectBegin()
	mock.
		ExpectExec("update reviews").
		WithArgs(review.UserId, review.ProductId, review.Rating, review.Text).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("update products").
		WithArgs(productRating, review.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = storage.UpdateReview(review, productRating)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// query error
	mock.ExpectBegin()
	mock.
		ExpectExec("update reviews").
		WithArgs(review.UserId, review.ProductId, review.Rating, review.Text).
		WillReturnError(fmt.Errorf("db_error"))
	mock.ExpectRollback()
	err = storage.UpdateReview(review, productRating)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func Test_DeleteReview(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	storage, err := NewReviewRepository(db, nil)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	review := models.Review{
		UserId:    1,
		ProductId: 1,
		Rating:    5,
		Text:      "Good",
	}

	var productRating float64

	// good query

	mock.ExpectBegin()
	mock.
		ExpectExec("delete from reviews").
		WithArgs(review.UserId, review.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.
		ExpectExec("update products").
		WithArgs(productRating, review.ProductId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = storage.DeleteReview(review.UserId, review.ProductId, productRating)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// query error
	mock.ExpectBegin()
	mock.
		ExpectExec("delete from reviews").
		WithArgs(review.UserId, review.ProductId).
		WillReturnError(fmt.Errorf("db_error"))
	mock.ExpectRollback()
	err = storage.DeleteReview(review.UserId, review.ProductId, productRating)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
