package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})

		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJson(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}
		log.Printf("failed to update task: %v", err)

		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка",
		})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
