package tests

import (
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

// Верни postgres Data Source Name из переменной окружения или значения по умолчанию.
func getDSN() string {
	env := os.Getenv("TODO_DSN")
	if env != "" {
		return env
	}
	// Поменяй user, password, dbname под свою конфигурацию!
	return "user=scheduler password=123456 dbname=scheduler sslmode=disable"
}

func openDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("postgres", getDSN())
	assert.NoError(t, err)
	return db
}

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, `SELECT count(id) FROM scheduler`)
}

func TestDB(t *testing.T) {
	db := openDB(t)
	defer db.Close()

	before, err := count(db)
	assert.NoError(t, err)

	today := time.Now().Format(`20060102`)

	var id int64
	err = db.QueryRowx(
		`INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES ($1, 'Todo', 'Комментарий', '') RETURNING id`,
		today,
	).Scan(&id)
	assert.NoError(t, err)

	var task Task
	err = db.Get(&task, `SELECT * FROM scheduler WHERE id=$1`, id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.ID)
	assert.Equal(t, `Todo`, task.Title)
	assert.Equal(t, `Комментарий`, task.Comment)

	_, err = db.Exec(`DELETE FROM scheduler WHERE id = $1`, id)
	assert.NoError(t, err)

	after, err := count(db)
	assert.NoError(t, err)

	assert.Equal(t, before, after)
}
