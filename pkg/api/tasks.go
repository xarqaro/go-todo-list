package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []*db.Task
	var err error

	search := r.URL.Query().Get("search")

	if search == "" {
		tasks, err = db.Tasks(r.Context(), 50)
		if err != nil {
			return
		}
	} else if t, err := time.Parse("02.01.2006", search); err == nil {
		formattedDate := t.Format(dateFormat)

		tasks, err = db.TasksByDate(r.Context(), formattedDate, 50)
		if err != nil {
			return
		}
	} else {
		tasks, err = db.TasksBySearch(r.Context(), search, 50)
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
