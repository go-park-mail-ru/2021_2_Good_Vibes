package postgresql

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB, err error) (*ReviewRepository, error) {
	if err != nil {
		return nil, err
	}

	return &ReviewRepository{
		db: db,
	}, nil
}

func (rb *ReviewRepository) GetAllRatingsOfProduct(productId int) ([]models.ProductRating, error) {
	rows, err := rb.db.Query("select rating, count(*) as count from reviews "+
		"where product_id=$1 "+
		"group by rating", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ratings []models.ProductRating
	for rows.Next() {
		rating := models.ProductRating{}
		err = rows.Scan(&rating.Rating, &rating.Count)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ratings, nil
}

func (rb *ReviewRepository) AddReview(review models.Review, productRating float64) error {
	err := tx(rb.db, func(tx *sql.Tx) error {
		_, err := rb.db.Exec(
			"insert into reviews(user_id, product_id, rating, text, date) values ($1, $2, $3, $4, $5)",
			review.UserId, review.ProductId, review.Rating, review.Text, review.Date)

		if err != nil {
			return err
		}

		_, err = rb.db.Exec("update products set rating=$1 where id=$2", productRating, review.ProductId)

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

func (rb *ReviewRepository) UpdateReview(review models.Review, productRating float64) error {
	err := tx(rb.db, func(tx *sql.Tx) error {
		_, err := rb.db.Exec(
			`update reviews set rating=$3, text=$4 where user_id=$1 and product_id=$2`,
			review.UserId, review.ProductId, review.Rating, review.Text)

		if err != nil {
			return err
		}

		_, err = rb.db.Exec("update products set rating=$1 where id=$2", productRating, review.ProductId)

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

func (rb *ReviewRepository) DeleteReview(userId int, productId int, productRating float64) error {
	err := tx(rb.db, func(tx *sql.Tx) error {
		_, err := rb.db.Exec(`delete from reviews where user_id=$1 and product_id=$2`,
			userId, productId)

		if err != nil {
			return err
		}

		_, err = rb.db.Exec("update products set rating=$1 where id=$2", productRating, productId)

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

func (rb *ReviewRepository) GetReviewsByProductId(productId int) ([]models.Review, error) {
	rows, err := rb.db.Query("select c.name, r.user_id, r.rating, r.text, r.date from reviews as r "+
		"join customers c on c.id = r.user_id "+
		"where r.product_id=$1 order by r.date desc", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.Review
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.UserName, &review.UserId, &review.Rating, &review.Text, &review.Date)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return reviews, nil
}

func (rb *ReviewRepository) GetReviewsByUser(userName string) ([]models.Review, error) {
	rows, err := rb.db.Query("select product_id, rating, text, date from reviews as r "+
		"join customers c on c.id = r.user_id "+
		"where c.name=$1 order by date", userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.Review
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.ProductId, &review.Rating, &review.Text, &review.Date)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return reviews, nil
}

func (rb *ReviewRepository) GetReviewByUserAndProduct(userId int, productId int) (models.Review, error) {
	var review models.Review
	row := rb.db.QueryRow("select user_id, product_id, rating, text from reviews where user_id=$1 and product_id=$2",
		userId, productId)
	err := row.Scan(&review.UserId, &review.ProductId, &review.Rating, &review.Text)

	if err == sql.ErrNoRows {
		return models.Review{}, nil
	}

	if err != nil {
		return models.Review{}, err
	}

	return review, nil
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
