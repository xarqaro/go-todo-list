package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"go1f/pkg/db"
	"go1f/pkg/server"
	"go1f/tests"
)

func run() error {
	var port string
	var dbFile string

	if port = os.Getenv("TODO_PORT"); port == "" {
		port = strconv.Itoa(tests.Port)
	}

	if dbFile = os.Getenv("TODO_DBFILE"); dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("init db: %w", err)
	}
	defer db.CloseDB()

	return server.Run(port)
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
