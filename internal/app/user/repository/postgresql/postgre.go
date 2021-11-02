package postgresql

import (
	"database/sql"
	customErrors "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/errors"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
	userHandler "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/product/delivery/http"
)

const customAvatar = userHandler.BucketUrl + "29654677-7947-46d9-a2e5-1ca33223e30d"

type StorageUserDB struct {
	db *sql.DB
}

func NewStorageUserDB(db *sql.DB, err error) (*StorageUserDB, error) {
	if err != nil {
		return nil, err
	}
	return &StorageUserDB{
		db: db,
	}, nil
}

func (su *StorageUserDB) GetUserDataByName(name string) (*models.UserDataStorage, error) {
	var tmp models.UserDataStorage
	row := su.db.QueryRow("select id, name, email, password from customers where name=$1", name)

	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func (su *StorageUserDB) InsertUser(newUser models.UserDataForReg) (int, error) {
	rows := su.db.QueryRow("insert into customers (name, email, password) values ($1, $2, $3) returning id",
		newUser.Name,
		newUser.Email,
		newUser.Password)

	var id int
	err := rows.Scan(&id)
	if err != nil {
		return customErrors.DB_ERROR, err
	}

	return id, nil
}

func (su *StorageUserDB) GetUserDataById(id uint64) (*models.UserDataStorage, error) {
	var tmp models.UserDataStorage
	row := su.db.QueryRow("select id, name, email, password, avatar from customers where id=$1", id)
	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password, &tmp.Avatar)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if !tmp.Avatar.Valid {
		tmp.Avatar.String = customAvatar
	}

	return &tmp, nil
}

func (su *StorageUserDB) SaveAvatarName(userId int, fileName string) error {
	_, err := su.db.Exec(`update customers set avatar = $2 where id = $1`, userId, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (su *StorageUserDB) UpdateUser(newData models.UserDataProfile) error {
	_, err := su.db.Exec(`update customers set name = $1, email = $2 where id = $3`, newData.Name,
		newData.Email, newData.Id)
	return err
}
