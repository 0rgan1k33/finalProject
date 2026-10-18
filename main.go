package main

import (
	"log"
	"net/http"

	"main.go/pkg/db"
)

const defPort = "7540"

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.Dir("./web")))
	if err := http.ListenAndServe(":"+defPort, nil); err != nil {
		log.Fatal(err)
	}
}
