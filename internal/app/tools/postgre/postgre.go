package postgre

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
)

func GetPostgres() (*sql.DB, error) {
	dsn := ""
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
