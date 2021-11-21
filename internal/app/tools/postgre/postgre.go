package postgre

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
)

func GetPostgres() (*sql.DB, error) {
	dsn := "postgres://dzuprfexsuwvev:cc5d3a25e89203423d7a86f094ccfefc8f6213021c5897cb39efca5a7cc429d4@ec2-52-19-96-181.eu-west-1.compute.amazonaws.com:5432/ddakq05jkvsdkh"
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
