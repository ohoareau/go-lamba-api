package handlers

import (
	"net/http"
)

//goland:noinspection GoUnusedParameter
func Robots(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("User-agent: *\nDisallow: /*"))

	if err != nil {
		http.Error(w, "Unable to write body", http.StatusInternalServerError)
		return
	}
}
