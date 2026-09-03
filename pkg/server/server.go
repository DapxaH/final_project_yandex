package server

import (
	"log"
	"net/http"
)

func Start() {
	const port = ":7540"
	const webDir = "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost%s", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
