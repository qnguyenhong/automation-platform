package httputil

import (
	"encoding/json"
	"net/http"

	apperr "github.com/qnguyenhong/automation-platform/pkg/errors"
)

// Envelope is a standard JSON response wrapper.
type Envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
	Meta  any    `json:"meta,omitempty"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data})
}

// JSONWithMeta writes a JSON response with data and metadata.
func JSONWithMeta(w http.ResponseWriter, status int, data any, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Data: data, Meta: meta})
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Error: message})
}

// AppError writes a JSON error response based on the AppError type.
func AppError(w http.ResponseWriter, err error) {
	var appErr *apperr.AppError
	if ok := errorAs(err, &appErr); ok {
		status := http.StatusInternalServerError
		switch appErr.Code {
		case "NOT_FOUND":
			status = http.StatusNotFound
		case "CONFLICT":
			status = http.StatusConflict
		case "UNAUTHORIZED":
			status = http.StatusUnauthorized
		case "FORBIDDEN":
			status = http.StatusForbidden
		case "BAD_REQUEST":
			status = http.StatusBadRequest
		}
		Error(w, status, appErr.Message)
		return
	}

	switch {
	case apperr.IsNotFound(err):
		Error(w, http.StatusNotFound, err.Error())
	case apperr.IsConflict(err):
		Error(w, http.StatusConflict, err.Error())
	case apperr.IsUnauthorized(err):
		Error(w, http.StatusUnauthorized, err.Error())
	case apperr.IsForbidden(err):
		Error(w, http.StatusForbidden, err.Error())
	case apperr.IsBadRequest(err):
		Error(w, http.StatusBadRequest, err.Error())
	default:
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}

// Decode reads the JSON body into the given value.
func Decode(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func errorAs(err error, target any) bool {
	type interfaceError interface {
		As(any) bool
	}
	if e, ok := err.(interfaceError); ok {
		return e.As(target)
	}
	return false
}
