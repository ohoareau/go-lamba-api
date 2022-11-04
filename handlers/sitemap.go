package handlers

import (
	"net/http"
)

//goland:noinspection GoUnusedParameter
func Sitemap(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<sitemapindex xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"></sitemapindex>"))

	if err != nil {
		http.Error(w, "Unable to write body", http.StatusInternalServerError)
		return
	}
}
