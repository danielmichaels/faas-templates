package function

import (
	"encoding/json"
	"handler/function/handlers"
	"log"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	pollen := handlers.ScrapePollenCount()
	err := handlers.Sender(pollen)
	if err != nil {
		log.Print("Error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pollen)
}
