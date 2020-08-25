package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/VamsiDevalla/garage_sale/internal/platform/web"
	"github.com/VamsiDevalla/garage_sale/internal/product"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	db  *sqlx.DB
	log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(p.db)
	if err != nil {
		return errors.Wrapf(err, "List handler :")
	}

	return web.Respond(w, list, http.StatusOK)
}

// Retrive gets a product from service layer and encodes it for the client respons.
func (p *Products) Retrive(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	prod, err := product.Retrive(p.db, id)

	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "Retrive Handler : getting product %q", id)
		}
	}

	return web.Respond(w, prod, http.StatusOK)
}

// Create decode the body of json to new product and saves it to the db
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct

	if err := web.Decode(r, &np); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	prod, err := product.Create(p.db, np, time.Now())

	if err != nil {
		return errors.Wrapf(err, "Create Handler :")
	}
	return web.Respond(w, prod, http.StatusOK)
}
