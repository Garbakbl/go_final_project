package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title" validate:"required"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4) RETURNING id;`
	err := db.QueryRow(query, task.Date, task.Title, task.Comment, task.Repeat).Scan(&id)
	return id, err
}

func Tasks(limit int, query string) ([]*Task, error) {
	var (
		tasks = make([]*Task, 0)
		err   error
		rows  *sql.Rows
	)
	if query == "" {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1", limit)
	} else if date, err := time.Parse("02.01.2006", query); err == nil {
		dateQuery := date.Format("20060102")
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 LIMIT $2", dateQuery, limit)
	} else {
		pattern := "%" + query + "%"
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $2 LIMIT $3", pattern, pattern, limit)
	}
	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if errors.Is(err, sql.ErrNoRows) {
			return tasks, nil
		}
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &task)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return &task, err
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET title = $1, comment = $2, repeat = $3 WHERE id = $4`
	res, err := db.Exec(query, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec("DELETE FROM scheduler WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf(`internal server error`)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf(`nothing to delete`)
	}
	return nil
}

func UpdateTaskDate(id, date string) error {
	_, err := db.Exec("UPDATE scheduler SET date = $1 WHERE id = $2", date, id)
	if err != nil {
		return err
	}
	return nil
}
