package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/tools/postgre"
	"strings"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB, err error) (*SearchRepository, error) {
	if err != nil {
		return nil, err
	}

	return &SearchRepository{
		db: db,
	}, nil
}

func (sr *SearchRepository) GetSuggests(str string) (models.Suggest, error) {
	var searchStr strings.Builder
	searchStr.WriteRune('%')
	searchStr.WriteString(str)
	searchStr.WriteRune('%')

	var products []models.ProductForSuggest
	var categories []models.CategoryForSuggest

	err := tx(sr.db, func(tx *sql.Tx) error {
		productRows, err := sr.db.Query(`select id, name, image from products where name ilike $1 order by  rating desc limit 5`, searchStr.String())
		if err != nil {
			return err
		}

		defer productRows.Close()

		var product models.ProductForSuggest

		for productRows.Next() {
			err := productRows.Scan(&product.Id, &product.Name, &product.Image)
			if err != nil {
				return err
			}

			products = append(products, product)
		}

		categoryRows, err := sr.db.Query(`select name, description from categories where description ilike $1 limit 5`, searchStr.String())
		if err != nil {
			return err
		}

		defer categoryRows.Close()

		var category models.CategoryForSuggest

		for categoryRows.Next() {
			err := categoryRows.Scan(&category.Name, &category.Description)
			if err != nil {
				return err
			}

			categories = append(categories, category)
		}

		return nil
	})

	if err != nil {
		return models.Suggest{}, err
	}

	suggests := models.Suggest{
		Products:   products,
		Categories: categories,
	}

	return suggests, nil
}

func (sr *SearchRepository) GetSearchResults(searchArray []string, filter postgre.Filter) ([][]models.Product, error) {
	var resultProducts [][]models.Product
	err := tx(sr.db, func(tx *sql.Tx) error {
		for _, str := range searchArray {
			products, err := sr.getSearchResultLocal(str, filter)
			if err != nil {
				return err
			}

			resultProducts = append(resultProducts, products)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resultProducts, nil
}

func (sr *SearchRepository) getSearchResultLocal(str string, filter postgre.Filter) ([]models.Product, error) {
	var searchStr strings.Builder
	searchStr.WriteRune('%')
	searchStr.WriteString(str)
	searchStr.WriteRune('%')

	if filter.OrderBy != postgre.TypeOrderRating && filter.OrderBy != postgre.TypeOrderPrice {
		filter.OrderBy = postgre.TypeOrderRating
	}
	if filter.TypeOrder != postgre.TypeOrderMin && filter.TypeOrder != postgre.TypeOrderMax {
		filter.TypeOrder = postgre.TypeOrderMin
	}
	filter.OrderBy = "p." + filter.OrderBy

	var products []models.Product

	/*rows, err := sr.db.Query(
	"select p.id, p.image, p.name, p.price, p.rating, c.name, " +
		"p.count_in_stock, p.description from products as p " +
		"join categories as c on c.id=p.category_id " +
		"where p.name ilike $1 and p.price >= $2 and p.price <= $3 " +
		"and p.rating >= $4 and p.rating <= $5 " +
		"order by "+filter.OrderBy+" "+filter.TypeOrder, searchStr.String(), filter.MinPrice,
	filter.MaxPrice, filter.MinRating, filter.MaxRating)*/

	rows, err := sr.db.Query("select p.id, p.image, p.name, p.price, p.rating, nc1.name, p.count_in_stock, p.description from products as p "+
		"join categories as nc1 on p.category_id = nc1.id "+
		"join categories as nc2 on nc1.lft >= nc2.lft AND "+
		"nc1.rgt <= nc2.rgt "+
		"where nc2.name = $1 and p.name ilike $6 "+
		"and p.price >= $2 and p.price <= $3 "+
		"and p.rating >= $4 and p.rating <= $5 "+
		"order by "+filter.OrderBy+" "+filter.TypeOrder, filter.NameCategory, filter.MinPrice,
		filter.MaxPrice, filter.MinRating, filter.MaxRating, searchStr.String())

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var product models.Product

	for rows.Next() {
		err := rows.Scan(&product.Id, &product.Image, &product.Name,
			&product.Price, &product.Rating, &product.Category,
			&product.CountInStock, &product.Description)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
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
