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
		                           "where product_id=$1 " +
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
			"insert into reviews(user_id, product_id, rating, text) values ($1, $2, $3, $4)",
			review.UserId,review.ProductId, review.Rating, review.Text)

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

func (rb *ReviewRepository) DeleteReview(userId int, productId int) error {
	_, err := rb.db.Exec(`delete from reviews where user_id=$1 and product_id=$2`,
		                        userId, productId)

	if err != nil {
		return err
	}

	return nil
}


func (rb *ReviewRepository) GetReviewsByProductId(productId int) ([]models.Review, error){
	rows, err := rb.db.Query("select c.name, r.rating, r.text from reviews as r " +
		                           "join customers c on c.id = r.user_id "+
		                           "where r.product_id=$1", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.Review
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.UserName, &review.Rating, &review.Text)
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
	rows, err := rb.db.Query("select product_id, rating, text from reviews as r " +
		                           "join customers c on c.id = r.user_id "+
	                          	   "where c.name=$1", userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.Review
	for rows.Next() {
		review := models.Review{}
		err = rows.Scan(&review.ProductId, &review.Rating, &review.Text)
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




/*
func (sb *BasketRepository) GetBasket(userId int) ([]models.BasketProduct, error) {
	var basketProducts []models.BasketProduct
	rows, err := sb.db.Query("select product_id, count from basket_products where user_id = $1 order by product_id", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		product := models.BasketProduct{}

		err := rows.Scan(&product.ProductId, &product.Number)
		if err != nil {
			return nil, err
		}

		basketProducts = append(basketProducts, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return basketProducts, nil
}

func (sb *BasketRepository) DropBasket(userId int) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`delete from basket where user_id=$1`, userId)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`delete from basket_products where user_id=$1`, userId)
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

func (sb *BasketRepository) DeleteProduct(product models.BasketProduct) error {
	err := tx(sb.db, func(tx *sql.Tx) error {
		_, err := tx.Exec(`delete from basket_products where user_id=$1 and product_id=$2`, product.UserId, product.ProductId)
		if err != nil {
			return err
		}

		rows, err := tx.Query(`select from basket_products where user_id=$1`, product.UserId)
		if err != nil {
			return err
		}

		defer rows.Close()
		if !rows.Next() {
			_, err := tx.Exec(`delete from basket where user_id=$1`, product.UserId)
			if err != nil {
				return err
			}
		}

		if rows.Err() != nil {
			return rows.Err()
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
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
*/
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