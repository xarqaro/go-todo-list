package api

import (
	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", NextDateHandler)
	r.Post("/api/task", AddTaskHandler)
	r.Get("/api/tasks", tasksHandler)
}
