package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VamsiDevalla/garage_sale/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Server is running on localhost:8080")

	defer log.Println("Server stopped")

	db, err := openDB()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Println("error applying migrations", err)
			os.Exit(1)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Println("error seeding database", err)
			os.Exit(1)
		}
		log.Println("Seed data complete")
		return
	}

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

	if _, err = w.Write(data); err != nil {
		log.Println("Error while writing to request", err)
	}
}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}
