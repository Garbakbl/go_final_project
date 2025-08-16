package db

import (
	"database/sql"
	"errors"
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

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int, query string) ([]*Task, error) {
	var (
		tasks = make([]*Task, 0)
		err   error
		rows  *sql.Rows
	)
	if query == "" {
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	} else if date, err := time.Parse("02.01.2006", query); err == nil {
		dateQuery := date.Format("20060102")
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?", dateQuery, limit)
	} else {
		pattern := "%" + query + "%"
		rows, err = db.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? LIMIT ?", pattern, pattern, limit)
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
	return tasks, nil
}
