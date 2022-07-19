package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

func (a *APIServer) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	fmt.Println("================================================================================================")
	a.logger.PrintError(trace)
	fmt.Println("================================================================================================")

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *APIServer) clientError(w http.ResponseWriter, status int, err error) {
	result := struct {
		Status string `json:"status"`
		Msg    string `json:"message"`
	}{
		Status: http.StatusText(status),
		Msg:    err.Error(),
	}

	bytes, errMarsh := json.Marshal(result)
	if errMarsh != nil {
		a.serverError(w, errMarsh)
		return
	}
	http.Error(w, string(bytes), status)
}

func (a *APIServer) notFound(w http.ResponseWriter) {
	_, err := io.WriteString(w, "this is my ")
	if err != nil {
		a.serverError(w, err)
	}
	a.clientError(w, http.StatusNotFound, errors.New("not found"))
}
