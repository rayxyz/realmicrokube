package main

import (
	"log"
	"net/http"
)

// Warning: This service is only for test.

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am in microtest!!!"))
	})
	log.Println("Microtest server running on port => 7777")
	http.ListenAndServe(":7777", nil)
}
