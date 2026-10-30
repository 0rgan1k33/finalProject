package main

import (
	"log"
	"net/http"

	"finalProject/pkg/api"
	"finalProject/pkg/db"
)

const defPort = "7540"

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Printf("Сервер запущен на http://localhost:%s/", defPort)
	if err := http.ListenAndServe(":"+defPort, nil); err != nil {
		log.Fatal(err)
	}
}
