package web

import (
	"encoding/json"
	"net/http"
)

// Respond will marshal the given data and sends it to the client with the given status code
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}

// RespondError can be used to respond to client in case of error
func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok{
		er := ErrorResponse{webErr.Err.Error()}
		if err := Respond(w, er, webErr.StatusCode); err != nil {
			return err
		}
		return nil
	}

	er := http.StatusText(http.StatusInternalServerError)
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
