package handlers

import (
	"log"
	"net/http"

	"github.com/VamsiDevalla/garage_sale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API return router for the app
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	productHandler := Products{
		db:  db,
		log: logger,
	}

	a := web.NewApp(logger)

	a.Handle(http.MethodGet, "/v1/products", productHandler.List)
	a.Handle(http.MethodGet, "/v1/products/{id}", productHandler.Retrive)
	a.Handle(http.MethodPost, "/v1/products", productHandler.Create)
	return a
}
