package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type RecommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB, err error) (*RecommendationRepository, error) {
	if err != nil {
		return nil, err
	}

	return &RecommendationRepository{
		db: db,
	}, nil
}

func (rr *RecommendationRepository) GetRecommendProductForUser(userId int) ([]models.ProductIdRecommendCount, error) {
	rows, err := rr.db.Query("SELECT product_id, counter FROM recommendation "+
		"where user_id = $1 order by counter DESC LIMIT 10", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productRecommend []models.ProductIdRecommendCount
	for rows.Next() {
		var currentProductRecommend models.ProductIdRecommendCount
		err = rows.Scan(&currentProductRecommend.Id, &currentProductRecommend.Counter)
		if err != nil {
			return nil, err
		}
		productRecommend = append(productRecommend, currentProductRecommend)
	}

	return productRecommend, nil
}


func (rr *RecommendationRepository) GetMostPopularProduct()([]models.Product, error) {
	rows, err := rr.db.Query("select id, image, name, price, rating, category_id, count_in_stock, description from products order by rating desc Limit 20")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		product := models.Product{}
		err = rows.Scan(&product.Id, &product.Image, &product.Name, &product.Price, &product.Rating, &product.Category, &product.CountInStock, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return products, nil
}
