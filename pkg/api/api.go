package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var Router chi.Router

func Init() {
	r := chi.NewRouter()

	r.Get("/api/nextdate", nextDayHandler)
	r.Post("/api/task", addTaskHandler)
	r.Get("/api/task", getTaskHandler)
	r.Put("/api/task", updateTaskHandler)
	r.Delete("/api/task", deleteTaskHandler)
	r.Get("/api/tasks", tasksHandler)
	r.Post("/api/task/done", doneTaskHandler)

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	Router = r
}
