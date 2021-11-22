package posgresql
//
//import (
//	"errors"
//	"fmt"
//	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
//	"gopkg.in/DATA-DOG/go-sqlmock.v1"
//	"reflect"
//	"testing"
//)
//
//func TestSelectAllCategories(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	storage, err := NewStorageCategoryDB(db, nil)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//
//	// good query
//	rows := sqlmock.
//		NewRows([]string{"nesting", "name"})
//
//	expect := []models.NestingCategory{
//		{0, "ALL_THINGS"},
//		{1, "CLOTHES"},
//		{2, "CLOTHES_MEN"},
//		{3, "CLOTHES_UP_MEN"},
//		{3, "SHOES_MEN"},
//	}
//
//	for _, item := range expect {
//		rows = rows.AddRow(item.Nesting, item.Name)
//	}
//
//	mock.
//		ExpectQuery("select ...").
//		WithArgs().
//		WillReturnRows(rows)
//
//	result, err := storage.SelectAllCategories()
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if !reflect.DeepEqual(result, expect) {
//		t.Errorf("results not match, want %v, have %v", expect, result)
//		return
//	}
//
//	// query error
//	mock.
//		ExpectQuery("select ...").
//		WithArgs().
//		WillReturnError(fmt.Errorf("db_error"))
//
//	_, err = storage.SelectAllCategories()
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
//	rows = sqlmock.NewRows([]string{"name"}).
//		AddRow("ALL_THINGS")
//
//	mock.
//		ExpectQuery("select ...").
//		WithArgs().
//		WillReturnRows(rows)
//
//	_, err = storage.SelectAllCategories()
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if err == nil {
//		t.Errorf("expected error, got nil")
//		return
//	}
//}
//
//func TestNewStorageCategoryDB_Fail(t *testing.T) {
//	db, _, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	err = errors.New("This is error: ")
//
//	_, err = NewStorageCategoryDB(db, err)
//
//	if err == nil {
//		t.Errorf("expected error, got nil")
//		return
//	}
//}
