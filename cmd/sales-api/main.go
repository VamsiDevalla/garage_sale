package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VamsiDevalla/garage_sale/internal/platform/database"
	"github.com/VamsiDevalla/garage_sale/cmd/sales-api/internal/handlers"
)

func main() {
	// =========================================================================
	// App Starting
	log.Println("Server is running on localhost:8080")
	defer log.Println("Server stopped")

	// =========================================================================
	// Start Database
	db, err := database.Open()
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	productHandler := handlers.Products{DB: db}

	// =========================================================================
	// Start API Service
	app := http.Server{
		Addr:         "localhost:8080",
		Handler:      http.HandlerFunc(productHandler.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)
	// Start the service listening for requests.
	go func() {
		serverErrors <- app.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
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
