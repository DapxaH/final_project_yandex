package db

import (
	"fmt"
)

// Task представляет задачу в планировщике
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в БД
func AddTask(task *Task) (int64, error) {
	query := `
			INSERT INTO scheduler (date, title, comment, repeat)
			VALUES (?, ?, ?, ?)
	`

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Tasks возвращает список задач, отсортированных по дате
func Tasks(limit int) ([]*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date
		LIMIT ?
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask возвращает задачу по её id
func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		WHERE id = ?	
	`

	var task Task

	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask обновляет поля существующей задачи
func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ? 
	`

	res, err := db.Exec(
		query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по id
func DeleteTask(id string) error {
	query := `
		DELETE FROM scheduler
		WHERE id = ?
	`

	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// UpdateDate изменяет дату следующего выполнения повторяющейся задачи (периодической)
func UpdateDate(next string, id string) error {
	query := `
		UPDATE scheduler
		SET date = ?
		WHERE id = ?
	`

	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")

	}

	return nil
}
