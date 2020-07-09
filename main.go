package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Server is running on localhost:8080")

	defer log.Println("Server stopped")

	app := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(listGames),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error: while serving %s", err)
	case <-shutdown:
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := app.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = app.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

type game struct {
	Name           string `json:"name"`
	PriceInDollars int    `json:"priceInDollars"`
	MyRating       int    `json:"myRating"`
}

func listGames(w http.ResponseWriter, r *http.Request) {
	gList := []game{}

	if true {
		gList = append(gList,
			game{"Fifa20", 60, 2},
			game{"Apex Legends", 0, 5})
	}

	data, err := json.Marshal(gList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while marshaling the game list", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _,err = w.Write(data); err != nil {
		log.Println("Error while writing to request", err)
	}
}
