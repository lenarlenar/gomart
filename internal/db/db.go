package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataBase struct {
	DB *sql.DB
}

func Open(dburi string) (*DataBase, error)   {
	db, err := sql.Open("postgres", dburi)

	if(err != nil) {
		return nil, err
	}
	return &DataBase {DB: db}, nil
}