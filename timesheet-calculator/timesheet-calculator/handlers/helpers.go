package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

// GetSecret retrieves the secret from openfaas and makes it available for use.
func GetSecret(secretName string) ([]byte, error) {
	secret, err := ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", secretName))
	if err != nil {
		return nil, err
	}

	s := strings.TrimSpace(string(secret))
	return []byte(s), nil
}
