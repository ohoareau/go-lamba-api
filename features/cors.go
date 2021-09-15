package features

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Cors(r *chi.Mux) {
	//goland:noinspection HttpUrlsUsage
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}
