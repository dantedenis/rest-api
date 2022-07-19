package apiserver

import (
	"net/http"
)

func (a *APIServer) middleware(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.PrintfInfo("Request Method: %s, URL: %s", r.Method, r.URL)
		h(w, r)
	}
}
