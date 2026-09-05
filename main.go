package main

import (
	"log"

	"github.com/DapxaH/final_project_yandex/pkg/db"
	"github.com/DapxaH/final_project_yandex/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server.Start()
}
