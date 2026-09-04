package api

import (
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

// TaskResp даёт JSON-ответ со списком задач.
type TaskResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// getTasksHandler возвращает список ближайших задач.
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
