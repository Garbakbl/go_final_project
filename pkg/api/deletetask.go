package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

// deleteTaskHandler удаляет задачу по id.
//
// @Summary      Удалить задачу
// @Description  Удаляет одну задачу по её id
// @Tags         tasks
// @Security     ApiKeyAuth
// @Param        id   query     string  true  "ID задачи"
// @Success      200  {object}  map[string]interface{}  "OK (пустой объект)"
// @Failure      404  {object}  map[string]interface{}  "Ничего не найдено"
// @Failure      500  {object}  map[string]interface{}  "Внутренняя ошибка"
// @Router       /api/task [delete]
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		switch err.Error() {
		case "internal server error":
			w.WriteHeader(http.StatusInternalServerError)
		case "nothing to delete":
			w.WriteHeader(http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{})
}
