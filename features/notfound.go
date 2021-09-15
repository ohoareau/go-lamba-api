package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/handlers"
)

func Notfound(r *chi.Mux) {
	r.NotFound(handlers.Notfound)
}
