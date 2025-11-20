package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect(conn string) (*sql.DB, error) {
	return sql.Open("postgres", conn)
}
