package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"os"
	"time"
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

	// Параметры ретрая
	const maxAttempts = 10
	const retryInterval = 2 * time.Second

	var err error
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			// Маловероятно, обычно ошибка будет на .Ping()
			time.Sleep(retryInterval)
			continue
		}

		err = db.Ping()
		if err == nil {
			break // успех, база доступна
		}

		// Соединиться не удалось — ждём и пробуем снова
		time.Sleep(retryInterval)
	}

	// Если не удалось подключиться после всех попыток, возвращаем ошибку
	if err != nil {
		return errors.New("could not connect to Postgres after retries: " + err.Error())
	}

	_, err = db.Exec(schema)
	return err
}
