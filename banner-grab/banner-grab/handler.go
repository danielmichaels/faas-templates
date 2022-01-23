package function

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"
)

type envelope map[string]interface{}

// Handle handles the request
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

	hostname, _ := os.Hostname()
	data := struct {
		Hostname string      `json:"hostname,omitempty"`
		IP       []string    `json:"ip,omitempty"`
		Headers  http.Header `json:"headers"`
		URL      string      `json:"url,omitempty"`
		Host     string      `json:"host,omitempty"`
		Method   string      `json:"method,omitempty"`
	}{
		Hostname: hostname,
		IP:       []string{},
		Headers:  r.Header,
		URL:      r.URL.RequestURI(),
		Host:     r.Host,
		Method:   r.Method,
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				data.IP = append(data.IP, ip.String())
			}
		}
	}

	err := writeJSON(w, http.StatusOK, envelope{"whoami": data}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
