package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"syscall"
	"context"
)

func main() {
	log.Println("Server is running on localhost:8080")

	defer log.Println("Server stopped")

	app := http.Server{
		Addr: "localhost:8080",
		Handler: http.HandlerFunc(echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout:  5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func(){
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

func echo(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3*time.Second);
	fmt.Fprintf(w, "you requested for %s %s", r.Method, r.URL.Path)
}
