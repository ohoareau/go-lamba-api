package handlers

import "net/http"

//goland:noinspection GoUnusedParameter
func Notfound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	_, err := w.Write([]byte(""))

	if err != nil {
		http.Error(w, "Unable to write body", http.StatusInternalServerError)
		return
	}
}
