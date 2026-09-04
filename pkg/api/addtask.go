package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}
