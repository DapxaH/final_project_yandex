package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJson(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}
		log.Printf("failed to delete task: %v", err)

		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка",
		})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
