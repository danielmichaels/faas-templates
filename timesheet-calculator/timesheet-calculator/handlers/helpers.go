package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error, message interface{}) {
	env := Envelope{"error": message}
	log.Println(err)
	errJ := WriteJSON(w, status, env, nil)
	if errJ != nil {
		log.Printf("error: %s", message)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return nil
}
