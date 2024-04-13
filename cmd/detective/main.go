package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	if err := http.ListenAndServe(":5714", nil); err != nil {
		log.Fatal(err)
	}
}
