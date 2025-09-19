package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type APIError struct {
	Error string `json:"error"`
}

// WriteJSONError writes a JSON error response with the given HTTP status code and error message.
// It automatically logs the error and sets the appropriate Content-Type header.
func WriteJSONError(w http.ResponseWriter, code int, msg string) {
	log.Printf("[API Error] %d: %s", code, msg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(APIError{Error: msg}) //TODO:: needs to return error

}

// ReadRequestBody reads and decodes the HTTP request body into the target struct.
// It limits the body size to 1MB and logs successful decoding for debugging.
// Returns an error if reading or JSON decoding fails.
func ReadRequestBody(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	log.Printf("[ReadRequestBody] Raw body: %s", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return err
	}
	log.Printf("[ReadRequestBody] Request body decoded successfully: %+v", target)
	return nil
}

// WriteJSONResponse writes a JSON response with the given HTTP status code and data.
// It automatically sets the Content-Type header and logs successful responses.
// Returns an error if JSON encoding fails.
func WriteJSONResponse(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	log.Printf("[WriteJSONResponse] Response written successfully: %d: %+v", code, data)
	return json.NewEncoder(w).Encode(data)
}
