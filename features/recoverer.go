package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Recoverer(r *chi.Mux) {
	r.Use(middleware.Recoverer)
}
