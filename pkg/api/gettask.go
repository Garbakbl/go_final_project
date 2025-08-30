package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

// getTaskHandler возвращает одну задачу по id.
//
// @Summary      Получить задачу
// @Description  Возвращает одну задачу по её id.
// @Tags         tasks
// @Security     ApiKeyAuth
// @Param        id   query     string  true  "ID задачи"
// @Success      200  {object}  db.Task  "Данные задачи"
// @Failure      404  {object}  map[string]interface{}  "Задача не найдена"
// @Router       /api/task [get]
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "задача не найдена"})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
