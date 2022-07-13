package apiserver

import (
	"net/http"
)

func (a *APIServer) middleware(h http.HandlerFunc) http.HandlerFunc {
	a.AddChan(1)
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		a.DoneChan()
	}
}
