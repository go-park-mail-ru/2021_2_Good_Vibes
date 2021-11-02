package posgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type NestedListCategory struct {
	Name          string
	LeftBoundary  int
	RightBoundary int
}

type StorageCategoryPostgres struct {
	db *sql.DB
}

func NewStorageCategoryDB(db *sql.DB, err error) (*StorageCategoryPostgres, error) {
	if err != nil {
		return nil, err
	}
	return &StorageCategoryPostgres{
		db: db,
	}, nil
}

func (sc *StorageCategoryPostgres) CreateCategory(categoryName string, parentCategoryName string) error {
	err := tx(sc.db, func(tx *sql.Tx) error {
		var parentCategory NestedListCategory
		row := tx.QueryRow(`select name, lft, rgt from categories where name=$1`, parentCategoryName)

		err := row.Scan(&parentCategory.Name, &parentCategory.LeftBoundary, &parentCategory.RightBoundary)
		if err != nil {
			return err
		}

		leftBoundary := parentCategory.RightBoundary
		rightBoundary := parentCategory.RightBoundary + 1

		var categoryId int

		err = tx.QueryRow(
			`insert into categories(name, lft, rgt) values ($1, $2, $3) returning id`,
			categoryName,
			leftBoundary,
			rightBoundary,
		).Scan(&categoryId)
		if err != nil {
			return err
		}

		_, err = tx.Exec("with a(id, name, lft, rgt) as "+
			"(select nc2.id, nc2.name as name, nc2.lft, nc2.rgt  from categories as nc1 "+
			"join categories nc2 on nc1.rgt < nc2.lft "+
			"where nc1.name = $1) "+
			"update categories set lft=lft+2,rgt=rgt+2 "+
			"where name in (select name from a)", parentCategoryName)

		if err != nil {
			return err
		}

		_, err = tx.Exec("with a(id, name, lft, rgt) as "+
			"(select nc2.id, nc2.name, nc2.lft, nc2.rgt from categories as nc1 "+
			"join categories nc2 on nc1.lft >= nc2.lft and nc1.rgt <= nc2.rgt "+
			"where nc1.name = $1) "+
			"update categories set rgt=rgt+2"+
			"where name in (select name from a)", parentCategoryName)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (sc *StorageCategoryPostgres) SelectAllCategories() ([]models.NestingCategory, error) {
	var nestingCategory []models.NestingCategory
	rows, err := sc.db.Query("select ((count(parent.name) - 1)::int), node.name as name " +
		"from categories as node, categories as parent " +
		"where node.lft between parent.lft and parent.rgt " +
		"group by node.name, node.lft " +
		"order by node.lft")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		category := models.NestingCategory{}

		err := rows.Scan(&category.Nesting, &category.Name)
		if err != nil {
			return nil, err
		}

		nestingCategory = append(nestingCategory, category)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return nestingCategory, nil
}

func tx(db *sql.DB, fb func(tx *sql.Tx) error) error {
	trx, _ := db.Begin()
	err := fb(trx)
	if err != nil {
		trx.Rollback()
		return err
	}
	trx.Commit()
	return nil
}
