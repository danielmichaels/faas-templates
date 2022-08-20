package function

import (
	"fmt"
	"handler/function/handlers"
	"net/http"
	"strings"
)

var routes = map[string]func(http.ResponseWriter, *http.Request){}

func init() {
	routes["/api"] = handlers.MakeAPIHandler()
	routes["/"] = handlers.MakeHomepageHandler()
	routes["/jwt"] = handlers.JwtExample()
	routes["/404"] = handlers.Make404Handler()
}

func Handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/api"):
		routes["/api"](w, r)
		return
	case strings.HasPrefix(r.URL.Path, "/jwt"):
		routes["/jwt"](w, r)
		return
	case fmt.Sprintf("%s", r.URL.Path) == "/":
		// must match explicitly, prefix of "/" will never not evaluate to true
		// meaning we'd never hit the default statement.
		routes["/"](w, r)
		return
	default:
		routes["/404"](w, r)
		return
	}
}
