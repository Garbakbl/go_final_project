package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var (
	schema = `CREATE TABLE IF NOT EXISTS scheduler (
    id SERIAL PRIMARY KEY,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(128) NOT NULL DEFAULT '',
    comment TEXT DEFAULT '',
    repeat VARCHAR(128) DEFAULT '');

	CREATE INDEX IF NOT EXISTS date_index ON scheduler (date);
    `

	db *sql.DB
)

func Init() error {
	// Получаем параметры из переменных окружения
	dsn := os.Getenv("PG_DSN")

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(schema)
	return err
}
