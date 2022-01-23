package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ipinfo/go/v2/ipinfo"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type envelope map[string]interface{}
type IPAddr struct {
	Ip string `json:"ip,omitempty"`
}
type WeatherApi struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   int     `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  int     `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
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

	wkey, err := getSecret("weather-key")
	if err != nil {
		log.Printf("secret failed to load %q", "weather-key")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ikey, err := getSecret("ipinfo-key")
	if err != nil {
		log.Printf("secret failed to load %q", "ipinfo-key")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	input := IPAddr{}
	err = readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ipInfo, err := getLocation(input, ikey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var weather WeatherApi
	err = getWeather(*ipInfo, string(wkey), &weather)

	err = writeJSON(w, http.StatusOK, envelope{"ipinfo": ipInfo, "weather": weather}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func getWeather(info ipinfo.Core, secretKey string, target *WeatherApi) error {
	apiUrl := `https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric`
	url := fmt.Sprintf(apiUrl, info.City, secretKey)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

func getLocation(ipaddr IPAddr, secretKey []byte) (*ipinfo.Core, error) {
	client := ipinfo.NewClient(nil, nil, string(secretKey))
	info, err := client.GetIPInfo(net.ParseIP(ipaddr.Ip))
	if err != nil {
		return nil, err
	}
	return info, nil
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
