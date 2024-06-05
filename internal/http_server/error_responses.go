package http_server

import (
	"encoding/json"
	"net/http"
)

func SendJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := map[string]string{"error_message": err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
