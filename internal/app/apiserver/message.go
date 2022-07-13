package apiserver

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

func (a *APIServer) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	a.logger.ErrorMsg.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *APIServer) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *APIServer) notFound(w http.ResponseWriter) {
	_, err := io.WriteString(w, "this is my ")
	if err != nil {
		a.logger.ErrorMsg.Println("Error writer:", err)
	}
	a.clientError(w, http.StatusNotFound)
}
