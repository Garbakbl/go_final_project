package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

// TasksResp — структура успешного ответа для списка задач.
// swagger:model TasksResp
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler возвращает список задач с опциональным фильтром поиска.
//
// @Summary      Получить список задач
// @Description  Возвращает список задач пользователя, с опциональным поиском по подстроке.
// @Tags         tasks
// @Security     ApiKeyAuth
// @Param        search  query     string  false  "Фильтр для поиска по задачам"
// @Success      200     {object}  api.TasksResp  "Список задач"
// @Failure      500     {object}  map[string]interface{} "Внутренняя ошибка"
// @Router       /api/tasks [get]
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	var (
		tasks []*db.Task
		err   error
	)
	query := r.FormValue("search")
	tasks, err = db.Tasks(50, query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&TasksResp{tasks})
}
