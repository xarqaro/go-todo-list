package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(255) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
); CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date)
`

func ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	if db == nil {
		return fmt.Errorf("No database to close.")
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}
	return nil
}

func Init(dbFile string) error {
	var install bool
	_, err := os.Stat(dbFile)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		install = true
	}
	ctx := context.Background()

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open db: %w", err)
	}

	if err := ping(ctx); err != nil {
		CloseDB()
		return fmt.Errorf("unable to connect to database: %w", err)
	}

	if install {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			CloseDB()
			return fmt.Errorf("transaction error: %w", err)
		}
		defer tx.Rollback()
		_, err = tx.Exec(schema)
		if err != nil {
			CloseDB()
			return err
		}
		if err := tx.Commit(); err != nil {
			CloseDB()
			return fmt.Errorf("failed to commit changes: %w", err)
		}
	}

	return nil
}
