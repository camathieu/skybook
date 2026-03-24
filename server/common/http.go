package common

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the standard error format for all API responses.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteError writes a JSON error response.
func WriteError(w http.ResponseWriter, message string, code int) {
	WriteJSON(w, ErrorResponse{Error: message, Code: code}, code)
}
