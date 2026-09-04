package api

import (
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, TaskResp{
		Tasks: tasks,
	})
}
