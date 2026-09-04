package api

import (
	"net/http"
	"time"

	"github.com/DapxaH/final_project_yandex/pkg/db"
)

// doneTaskHandler отмечает задачу выполненной (пустой кружочек --> галочка)
// Одноразовая задача удаляется, а периодическая переносится на следующую дату.
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJson(w, map[string]string{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateDate(next, id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJson(w, map[string]string{})
}
