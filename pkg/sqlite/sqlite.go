package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New(dbFile string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return sqlDB, nil
}
