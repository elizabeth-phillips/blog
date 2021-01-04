package psql

import (
	"database/sql"
	"fmt"
)

//NewPostgresConnection creates a connection to a Postgres database
func NewPostgresConnection(connection string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Ping successful")

	return
}
