package api

import (
	"encoding/json"
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})

		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, map[string]string{})
}
