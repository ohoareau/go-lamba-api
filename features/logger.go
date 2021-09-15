package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Logger(r *chi.Mux) {
	r.Use(middleware.Logger)
}
