package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const dbPath = "../../internal/database/tracker.db"

// Open открывает соединение к БД sqlite tracker.db
func Open() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite", dbPath)

	return
}
