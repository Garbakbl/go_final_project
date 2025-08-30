package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

type Response struct {
	ID int64 `json:"id"`
}

func checkDate(task *db.Task) error {
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if t.Before(now.AddDate(0, 0, -1)) {
		if task.Repeat == "" {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}
	return nil
}

// @Summary      Добавить задачу
// @Description  Создаёт новую задачу
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body db.Task true "Новая задача"
// @Success      200 {object} api.Response
// @Failure      400 {object} map[string]any "Ошибка клиента"
// @Failure      500 {object} map[string]any "Ошибка сервера"
// @Router       /api/task [post]
// @Security     ApiKeyAuth
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "decoder error"})
		return
	}

	err = validator.New().Struct(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "validator error"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": "date error"})
		return
	}

	taskID, err := db.AddTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": "db error" + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{ID: taskID})
}
