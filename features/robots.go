package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/handlers"
)

func Robots(r *chi.Mux) {
	r.Get("/robots.txt", handlers.Robots)
}
