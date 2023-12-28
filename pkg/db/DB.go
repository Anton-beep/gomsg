package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", user, password, dbname)
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
