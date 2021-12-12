package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

type StorageNotifyPostgres struct {
	db *sql.DB
}

func NewStorageNotifyDB(db *sql.DB, err error) (*StorageNotifyPostgres, error) {
	if err != nil {
		return nil, err
	}
	return &StorageNotifyPostgres{
		db: db,
	}, nil
}


func (sn *StorageNotifyPostgres) GetStatusChanges() ([]models.ChangedStatus, error) {
	var changedStatuses []models.ChangedStatus
	rows, err := sn.db.Query(`select id, user_id, status, email from orders where status_meta='changed'`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		changedStatus := models.ChangedStatus{}

		err := rows.Scan(&changedStatus.OrderId, &changedStatus.UserId, &changedStatus.Status, &changedStatus.Email)
		if err != nil {
			return nil, err
		}

		changedStatuses = append(changedStatuses, changedStatus)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return changedStatuses, nil
}

func (sn *StorageNotifyPostgres) GetAddressInfo(orderId int) (models.Address, error) {
	address := models.Address{}

	err := sn.db.QueryRow("select country, region, city, street, house, flat, a_index from delivery_address where order_id = $1", orderId).
		Scan(
			&address.Country,
			&address.Region,
			&address.City,
			&address.Street,
			&address.House,
			&address.Flat,
			&address.Index,
		)

	if err == sql.ErrNoRows {
		return models.Address{}, nil
	}
	if err != nil {
		return models.Address{}, err
	}
	return address, nil
}