package api

import (
	"log"
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
		log.Printf("failed to get tasks: %v", err)

		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка сервера",
		})
		return
	}

	writeJson(w, http.StatusOK, TaskResp{
		Tasks: tasks,
	})
}
