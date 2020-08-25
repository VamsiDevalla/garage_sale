package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

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
func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
