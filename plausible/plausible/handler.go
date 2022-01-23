package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type envelope map[string]interface{}
type PlausibleAggregate struct {
	Results struct {
		BounceRate struct {
			Value int `json:"value"`
		} `json:"bounce_rate"`
		Pageviews struct {
			Value int `json:"value"`
		} `json:"pageviews"`
		VisitDuration struct {
			Value int `json:"value"`
		} `json:"visit_duration"`
		Visitors struct {
			Value int `json:"value"`
		} `json:"visitors"`
	} `json:"results"`
}
type SiteData struct {
	SiteId    string `json:"site_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	ApiKey    string `json:"-"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	validateCORS(w, r)

	apiKey, err := getSecret("plausible-key")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	input := SiteData{}
	err = readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	input.ApiKey = strings.TrimSpace(string(apiKey))

	analytics := PlausibleAggregate{}
	err = plausibleAggregate(input, &analytics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("err: %#v", err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"analytics": analytics}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkRedirectFunc(r *http.Request, via []*http.Request) error {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", via[0].Header.Get("Authorization")))
	return nil
}
func plausibleAggregate(input SiteData, result *PlausibleAggregate) error {
	// take in date from client in format 2021-12-01,2021-12-31
	apiUrl := "https://plausible.io/api/v1/stats/aggregate?site_id=%s&period=custom&date=%s,%s&metrics=visitors,pageviews,bounce_rate,visit_duration"
	url := fmt.Sprintf(apiUrl, input.SiteId, input.StartDate, input.EndDate)
	log.Println(url)

	client := &http.Client{Timeout: 10 * time.Second}
	client.CheckRedirect = checkRedirectFunc
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer: %s", input.ApiKey))
	log.Println(req.Header.Get("Authorization"))

	log.Printf("request: %#v", req)
	r, err := client.Do(req)
	if err != nil {
		log.Println("error do", err)
		return err
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.Body)
		return errors.New(fmt.Sprintf("unexpected status: %d", r.StatusCode))
	}
	defer r.Body.Close()

	log.Printf("body: %#v", r.Body)

	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return err
	}
	return nil
}

func validateCORS(w http.ResponseWriter, r *http.Request) {
	origins := strings.Split(os.Getenv("origins"), ",")

	if r.Method == "OPTIONS" {
		for _, origin := range origins {
			if r.Header.Get("Origin") == origin {
				w.Header().Set("Access-Control-Allow-Headers", "Authorization")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				w.Header().Add("Access-Control-Allow-Origin", origin)
				w.Header().Add("Access-Control-Max-Age", "300")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
	}

	for _, origin := range origins {
		if r.Header.Get("Origin") == origin {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}
	}
	return
}

func getSecret(secretName string) ([]byte, error) {
	secret, err := ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", secretName))
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(js)
	return nil
}

// readJSON is helper for trapping errors and return values for JSON related
// handlers
func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Set a max body length. Without this it will accept unlimited size requests
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Init a Decoder and call DisallowUnknownFields() on it before decoding.
	// This means that JSON from the client will be rejected if it contains keys
	// which do not match the target destination struct. If not implemented,
	// the decoder will silently drop unknown fields - this will raise an error instead.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode the request body into the target struct/destination
	err := dec.Decode(dst)
	if err != nil {
		// start triaging the various JSON related errors
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshallError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the
		// *json.SyntaxError. If it does, then return a user-readable error
		// message including the location of the problem
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Decode() can also return an io.ErrUnexpectedEOF for JSON syntax errors. This is
		// checked for with errors.Is() and returns a generic error message to the client.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Wrong JSON types will return an error when they do not match the target destination
		// struct.
		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshallTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Offset)

		// An EOF error will be returned by Decode() if the request body is empty. Use errors.Is()
		// to check for this and return a human-readable error message
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// If JSON contains a field which cannot be mapped to the target destination
		// then Decode will return an error message in the format "json: unknown field "<name>""
		// We check for this, extract the field name and interpolate it into an error
		// which is returned to the client
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		// If the request body exceeds maxBytes the decode will fail with a
		// "http: request body too large".
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// A json.InvalidUnmarshallError will be returned if we pass a non-nil pointer
		// to Decode(). We catch and panic, rather than return an error.
		case errors.As(err, &invalidUnmarshallError):
			panic(err)

		// All else fails, return an error as-is
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to anonymous empty struct as the
	// destination. If the body only has one JSON value then an io.EOF error
	// will be returned. If there is anything else, extra data has been sent
	// and we craft a custom error message back to the client
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
