package postgresql

import (
	"database/sql"
	models "github.com/go-park-mail-ru/2021_2_Good_Vibes/internal/app/models"
)

const customAvatar = "https://products-bucket-ozon-good-vibes.s3.eu-west-1.amazonaws.com/29654677-7947-46d9-a2e5-1ca33223e30d"

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
	row := su.db.QueryRow("SELECT id, name, email, password, birthday, real_name, real_surname, sex FROM customers WHERE name=$1", name)

	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password, &tmp.BirthDay, &tmp.RealName, &tmp.RealSurname, &tmp.Sex)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func (su *StorageUserDB) GetUserDataById(id uint64) (*models.UserDataStorage, error) {
	var tmp models.UserDataStorage
	row := su.db.QueryRow("SELECT id, name, email, password, avatar, birthday, real_name, real_surname, sex FROM customers WHERE id=$1", id)
	err := row.Scan(&tmp.Id, &tmp.Name, &tmp.Email, &tmp.Password, &tmp.Avatar, &tmp.BirthDay, &tmp.RealName, &tmp.RealSurname, &tmp.Sex)

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
	_, err := su.db.Exec(`UPDATE customers SET avatar = $2 WHERE id = $1`, userId, fileName)
	if err != nil {
		return err
	}
	return nil
}

func (su *StorageUserDB) UpdateUser(newData models.UserDataProfile) error {
	_, err := su.db.Exec(`UPDATE customers SET name = $1, email = $2, birthday = $4, real_name = $5, real_surname = $6, sex = $7  WHERE id = $3`, newData.Name,
		newData.Email, newData.Id, newData.BirthDay, newData.RealName, newData.RealSurname, newData.Sex)

	if err != nil {
		return err
	}
	return nil
}

func (su *StorageUserDB) UpdatePassword(newData models.UserDataPassword) error {
	_, err := su.db.Exec(`UPDATE customers SET password = $1 WHERE id = $2`, newData.Password,
		newData.Id)

	if err != nil {
		return err
	}
	return nil
}

func (su *StorageUserDB) InsertUserToken(email string, token string) error {
	_, err := su.db.Exec(`INSERT INTO email_confirm (email, status, token) VALUES ($1, 0, $2)`, email, token)

	if err != nil {
		return err
	}

	return nil
}
