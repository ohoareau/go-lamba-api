package main

import (
	"github.com/ohoareau/gola"
	"github.com/ohoareau/gola/common"
	"net/http"
)

func main() {
	gola.Main(common.Options{
		Apigw2Configurator: configure,
		Features:           featurize(),
	})
}

func configure(r common.HttpRouter) {
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("Hello!"))
	})
}

func featurize() common.Features {
	return common.Features{
		"logger":    true,
		"recoverer": true,
		"cors":      true,
		"ping":      true,
		"root":      true,
		"profiler":  true,
		"notfound":  true,
	}
}
