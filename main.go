package main

import (
	"log"
	"os"
	"strconv"
	"todo-list/pkg/server"
	"todo-list/tests"
)

func main() {
	var port string

	if port = os.Getenv("TODO_PORT"); port == "" {
		port = strconv.Itoa(tests.Port)
	}

	if err := server.Run(port); err != nil {
		log.Fatal(err)
	}
}
