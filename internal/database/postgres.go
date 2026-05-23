package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Connect(
	host string,
	port int,
	user string,
	password string,
	dbname string) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	db, err := sql.Open(
		"postgres",
		dsn,
	)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
