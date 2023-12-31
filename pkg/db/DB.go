package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	user     = "admin"
	password = "admin"
	dbname   = "postgres"
)

type APIDB struct {
	db *sql.DB
}

func NewDb() (*APIDB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		"localhost", // change this to ip of local network for docker container
		5432,
		dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &APIDB{db}, nil
}

// checks if user deleted something or not
func handleResultAfterEdit(result sql.Result) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		if err := rows.Close(); err != nil {
			zap.L().Error(err.Error())
		}
	}
}
