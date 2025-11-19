package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aliharirian/TerraPeak/logger"
	"github.com/go-chi/chi/v5/middleware"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	TotalPages int  `json:"totalPages"`
	TotalItems int  `json:"totalItems"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

// SuccessResponse creates a successful API response
func SuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta: &MetaInfo{
			Version:   "v1",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: middleware.GetReqID(r.Context()),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode success response: %v", err)
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(w http.ResponseWriter, r *http.Request, code, message string, statusCode int) {
	response := APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Meta: &MetaInfo{
			Version:   "v1",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			RequestID: middleware.GetReqID(r.Context()),
		},
	}

	logger.Infof("API error response: %s (requestId: %s, code: %s, status: %d, path: %s)",
		message, middleware.GetReqID(r.Context()), code, statusCode, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode error response: %v", err)
	}
}

// NotFoundResponse creates a 404 error response
func NotFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "The requested resource was not found"
	}
	ErrorResponse(w, r, "RESOURCE_NOT_FOUND", message, http.StatusNotFound)
}

// BadRequestResponse creates a 400 error response
func BadRequestResponse(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Invalid request parameters"
	}
	ErrorResponse(w, r, "BAD_REQUEST", message, http.StatusBadRequest)
}

// InternalErrorResponse creates a 500 error response
func InternalErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "An internal server error occurred"
	}
	ErrorResponse(w, r, "INTERNAL_ERROR", message, http.StatusInternalServerError)
}

// ValidationErrorResponse creates a 422 error response
func ValidationErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Validation failed"
	}
	ErrorResponse(w, r, "VALIDATION_ERROR", message, http.StatusUnprocessableEntity)
}
