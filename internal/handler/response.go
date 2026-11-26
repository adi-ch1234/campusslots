package handler
import (
	"encoding/json"
	"net/http"
)
// ErrorResponse is the standard JSON shape for every API error.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
// WriteJSON encodes payload as JSON and writes it with the given HTTP status code.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
// WriteError sends a consistent machine-readable error code and human message.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}