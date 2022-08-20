package handlers

import (
	"net/http"
	"os"
)

func MakeHomepageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./static/index.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Write(data)
	}
}
