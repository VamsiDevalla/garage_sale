package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List gets all Products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created FROM products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Retrive gets single Product with matching id
func Retrive(db *sqlx.DB, id string) (*Product, error) {
	product := Product{}
	const q = `SELECT product_id, name, cost, quantity, date_updated, date_created FROM products WHERE product_id=$1`

	if err := db.Get(&product, q, id); err != nil {
		return nil, errors.Wrapf(err, "selecting a single product ", id)
	}

	return &product, nil
}

// Create saves a new product to the db
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}

	const q = `INSERT INTO products
	(product_id, name, cost, quantity, date_created, date_updated)
	VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "inserting into db: %v", p)
	}

	return &p, nil
}
