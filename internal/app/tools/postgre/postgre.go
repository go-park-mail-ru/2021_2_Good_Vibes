package postgre

import (
	"database/sql"
	"fmt"
	configApp "github.com/go-park-mail-ru/2021_2_Good_Vibes/config"
	_ "github.com/jackc/pgx/stdlib"
)

func GetPostgres() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		configApp.ConfigApp.DataBase.User, configApp.ConfigApp.DataBase.DBName,
		configApp.ConfigApp.DataBase.Password, configApp.ConfigApp.DataBase.Host,
		configApp.ConfigApp.DataBase.Port)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	return db, nil
}
