package handlers

import (
	"net/http"
)

//goland:noinspection GoUnusedParameter
func Root(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("{}"))

	if err != nil {
		http.Error(w, "Unable to write body", http.StatusInternalServerError)
		return
	}
}
