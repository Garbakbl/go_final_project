package db

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

var (
	schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(128) DEFAULT "");

	CREATE INDEX date_index ON scheduler (date);
    `

	DBFilePath = os.Getenv("TODO_DBFILE")

	db *sql.DB
)

func Init(dbFile string) error {
	var newDB bool
	_, err := os.Stat(dbFile)
	if err != nil {
		newDB = true
		err = nil
	}

	db, err = sql.Open("sqlite", dbFile)
	err = db.Ping()

	if newDB {
		_, err = db.Exec(schema)
	}
	return err
}
