package server

import (
	"log"
	"net/http"

	"github.com/DapxaH/final_project_yandex/pkg/api"
)

func Start() {
	const port = ":7540"
	const webDir = "./web"

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost%s", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
