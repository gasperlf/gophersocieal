package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return json.NewEncoder(w).Encode(struct{}{})
	}

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) error {

	w.Header().Set("Content-Type", "application/json")

	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func noContentResponse(w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json")
	return writeJSON(w, status, nil)
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}
