package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"go_final_project/pkg/db"
	"net/http"
)

// updateTaskHandler обновляет существующую задачу.
//
// @Summary      Обновить задачу
// @Description  Обновляет все поля существующей задачи по её id. Все поля задачи должны быть заполнены.
// @Tags         tasks
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        task  body      db.Task  true  "Данные задачи для обновления"
// @Success      200   {object}  map[string]interface{}  "Успешно (пустой объект)"
// @Failure      400   {object}  map[string]interface{}  "Ошибка валидации/даты/декодирования"
// @Failure      404   {object}  map[string]interface{}  "Задача не найдена"
// @Router       /api/task [put]
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.UpdateTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"error": "задача не найдена"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{})
}
