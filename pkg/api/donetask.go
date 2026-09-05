package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

// doneTaskHandler отмечает задачу выполненной (пустой кружочек --> галочка)
// Одноразовая задача удаляется, а периодическая переносится на следующую дату.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJson(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}
		log.Printf("failed to get task: %v", err)

		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка",
		})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			if errors.Is(err, db.ErrTaskNotFound) {
				writeJson(w, http.StatusNotFound, map[string]string{
					"error": "Задача не найдена",
				})
				return
			}
			log.Printf("failed to delete completed task: %v", err)

			writeJson(w, http.StatusInternalServerError, map[string]string{
				"error": "Внутренняя ошибка",
			})
			return
		}

		writeJson(w, http.StatusOK, map[string]string{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateDate(next, id)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJson(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
			return
		}
		log.Printf("failed to update task date: %v", err)

		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "Внутренняя ошибка",
		})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}
