package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var Router chi.Router

func Init() {
	r := chi.NewRouter()

	r.Post("/api/signin", signin)

	r.Get("/api/nextdate", nextDayHandler)
	r.With(auth).Post("/api/task", addTaskHandler)
	r.With(auth).Get("/api/task", getTaskHandler)
	r.With(auth).Put("/api/task", updateTaskHandler)
	r.With(auth).Delete("/api/task", deleteTaskHandler)
	r.With(auth).Get("/api/tasks", tasksHandler)
	r.With(auth).Post("/api/task/done", doneTaskHandler)

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	Router = r
}
