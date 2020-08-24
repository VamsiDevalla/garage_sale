package handlers

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

var sm = http.NewServeMux()

// API return router for the app
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	productHandler := Products{
		db:  db,
		log: logger,
	}

	sm.HandleFunc("/v1/products", productHandler.List)
	sm.HandleFunc("/v1/products/a2b0639f-2cc6-44b8-b97b-15d69dbb511e", productHandler.Retrive)
	sm.HandleFunc("/v1/products/create", productHandler.Create)
	return sm
}
