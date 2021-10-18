package posgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

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
	return nestingCategory, nil
}
