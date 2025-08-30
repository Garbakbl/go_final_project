package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

// doneTaskHandler отмечает задачу как выполненную (с учётом повторяемости).
//
// @Summary      Отметить задачу как выполненную
// @Description  Отмечает выполнение задачи. Если задача одиночная — удаляет её. Если повтаряющаяся — переносит на следующую дату.
// @Tags         tasks
// @Security     ApiKeyAuth
// @Param        id   query     string  true  "ID задачи"
// @Success      200  {object}  map[string]interface{}  "Выполнено успешно (пустой объект)"
// @Failure      404  {object}  map[string]interface{}  "Задача не найдена"
// @Failure      500  {object}  map[string]interface{}  "Внутренняя ошибка"
// @Router       /api/task/done [post]
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "задача не найдена"})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{})
		return
	}
	newDate, _ := NextDate(time.Now(), task.Date, task.Repeat)
	err = db.UpdateTaskDate(task.ID, newDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{})
}
