package handlers

import (
	"encoding/json"
	"net/http"
)

type UserProfile struct {
	Homepage string `json:"homepage"`
	Twitter  string `json:"twitter"`
	GitHub   string `json:"github"`
	Gumroad  string `json:"gumroad"`
}

var user UserProfile

func init() {
	user = UserProfile{
		Homepage: "https://alexelis.io",
		Twitter:  "https://twitter.com/alexelisuk",
		GitHub:   "https://github.com/alexellis/",
		Gumroad:  "https://store.openfaas.com/",
	}
}

func MakeAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(data)
	}
}
