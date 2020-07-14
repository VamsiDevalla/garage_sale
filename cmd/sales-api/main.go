package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VamsiDevalla/garage_sale/cmd/sales-api/internal/handlers"
	"github.com/VamsiDevalla/garage_sale/internal/platform/conf"
	"github.com/VamsiDevalla/garage_sale/internal/platform/database"
)

func main() {

	// =========================================================================
	// Configuration
	// export SALES_DB_DISABLE_SSL=true to export DisableSSL to env properties

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Path       string `conf:"default:postgres"`
			DisableSSL bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			log.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	// =========================================================================
	// App Starting
	log.Println("Server is running on", cfg.Web.Address)
	defer log.Println("Server stopped")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error : generating config for output : %v", err)
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Path:       cfg.DB.Path,
		DisableSSL: cfg.DB.DisableSSL,
	})
	if err != nil {
		log.Fatalf("error: connecting to db: %s", err)
	}
	defer db.Close()

	productHandler := handlers.Products{DB: db}

	// =========================================================================
	// Start API Service
	app := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(productHandler.List),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := app.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = app.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}
