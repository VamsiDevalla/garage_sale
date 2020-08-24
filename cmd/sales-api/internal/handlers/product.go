package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/VamsiDevalla/garage_sale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) List(w http.ResponseWriter, r *http.Request) {
	list, err := product.List(p.db)
	if err != nil {
		p.log.Printf("error: listing products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		p.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.log.Println("error writing result", err)
	}
}

// Retrive gets a product from service layer and encodes it for the client respons.
func (p *Products) Retrive(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := product.Retrive(p.db, id)

	if err != nil {
		p.log.Printf("error: listing products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		p.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.log.Println("error writing result", err)
	}
}

// Create decode the body of json to new product and saves it to the db
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {
	var np product.NewProduct
	json.NewDecoder(r.Body).Decode(&np)

	product, err := product.Create(p.db, np, time.Now())

	if err != nil {
		p.log.Println("creating product", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		p.log.Println("error marshalling result", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		p.log.Println("error writing result", err)
	}
}
