package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/ohoareau/gola/handlers"
)

func Sitemap(r *chi.Mux) {
	r.Get("/sitemap.xml", handlers.Sitemap)
}
