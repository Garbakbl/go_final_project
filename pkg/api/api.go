// @title           Tasks API
// @version         1.0
// @description     ...
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/swaggo/http-swagger"
	_ "go_final_project/docs"
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	_, thisFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(thisFile), "..", "..") // путь к корню проекта
	absWeb := filepath.Join(projectRoot, "web")
	r.Handle("/*", http.FileServer(http.Dir(absWeb)))

	Router = r
}
