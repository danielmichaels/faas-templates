package function

import (
	"log"
	"net/http"
	"os"
)

var routes = map[string]func(http.ResponseWriter, *http.Request){}

func init() {
	routes["/"] = MakeHomepageHandler()
}

func isDockerCompose() string {
	file := os.Getenv("DOCKER_COMPOSE")
	index := "./function/static/index.html"
	if file == "" {
		index = "./static/index.html"
	}
	return index
}

func MakeHomepageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index := isDockerCompose()
		data, err := os.ReadFile(index)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			log.Println(err)
			return
		}

		w.Write(data)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {
	routes["/"](w, r)
	return
}
