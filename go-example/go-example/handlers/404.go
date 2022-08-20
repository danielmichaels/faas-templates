package handlers

import (
	"net/http"
	"os"
)

func Make404Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./static/404.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Write(data)
	}
}
