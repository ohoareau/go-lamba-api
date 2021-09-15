package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"os"
)

func Profiler(r *chi.Mux) {
	if 0 != len(os.Getenv("DEBUG")) {
		r.Mount("/debug", middleware.Profiler())
	}
}
