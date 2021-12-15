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

func (rr *RecommendationRepository) GetRecommendProductForUser (userId int) ([]models.ProductIdRecommendCount, error) {
	rows, err := rr.db.Query("SELECT product_id, counter FROM recommendation " +
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
