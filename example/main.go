package main

import (
	"github.com/ohoareau/gola"
	"net/http"
)

func main() {
	gola.Gola(configure, featurize)
}

func configure(r gola.Router) {
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("Hello!"))
	})
}

func featurize() gola.Features {
	return gola.Features{
		"logger":    true,
		"recoverer": true,
		"cors":      true,
		"ping":      true,
		"root":      true,
		"profiler":  true,
		"notfound":  true,
	}
}
