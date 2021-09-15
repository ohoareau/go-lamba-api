package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/handlers"
)

func Root(r *chi.Mux) {
	r.Get("/", handlers.Root)
}
