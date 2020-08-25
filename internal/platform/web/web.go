package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Handler is the signature used by all handlers in the service
type Handler func(w http.ResponseWriter, r *http.Request) error

// App handels all the multiplexing for the api
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp creates an new app and return pointer to it
func NewApp(log *log.Logger) *App {
	return &App{
		mux: chi.NewMux(),
		log: log,
	}
}

// Handle attaches a http HandlerFunc method to the string pattern and method
func (a *App) Handle(method, pattern string, h Handler) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := h(w,r)

		if err != nil {

			// Log the error.
			a.log.Printf("ERROR : %+v", err)

			// Respond to the error.
			if err := RespondError(w, err); err != nil {
				a.log.Printf("ERROR : %v", err)
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
