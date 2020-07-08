package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	h := http.HandlerFunc(echo)

	log.Println("Server is running on localhost:8080")

	if err := http.ListenAndServe("localhost:8080", h); err != nil {
		log.Fatal("Unable to start server", err)
	}

}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you requested for %s %s", r.Method, r.URL.Path)
}
