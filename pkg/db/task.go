package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const query = "SELECT id, date, title, comment, repeat FROM scheduler"
const selectWithOrder = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
const selectBySearch = query + " WHERE LOWER(title) LIKE ? OR LOWER(comment) LIKE ? ORDER BY date LIMIT ?"
const selectByDate = query + " WHERE date = ? ORDER BY date LIMIT ?"

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	tasks := []*Task{}

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func Tasks(ctx context.Context, limit int) ([]*Task, error) {
	tasks := []*Task{}

	rows, err := db.QueryContext(ctx, selectWithOrder, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err = scanTasks(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func TasksBySearch(ctx context.Context, search string, limit int) ([]*Task, error) {
	tasks := []*Task{}
	pattern := "%" + search + "%"

	rows, err := db.QueryContext(ctx, selectBySearch, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err = scanTasks(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func TasksByDate(ctx context.Context, date string, limit int) ([]*Task, error) {
	fmt.Println("tasksByDate:")
	fmt.Println("date", date)
	tasks := []*Task{}

	rows, err := db.QueryContext(ctx, selectByDate, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks, err = scanTasks(rows)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func AddTask(ctx context.Context, task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.ExecContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return id, nil
}
