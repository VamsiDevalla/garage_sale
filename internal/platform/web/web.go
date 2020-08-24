package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type App struct {
	mux *chi.Mux
	log *log.Logger
}

func NewApp(log *log.Logger) *App {
	return &App{
		mux: chi.NewMux(),
		log: log,
	}
}

func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
