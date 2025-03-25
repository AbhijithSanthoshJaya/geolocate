// ErrorResponse struct for API errors
package client

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Status int `json:"status"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("HttpError: %d , %s", e.Status, getErrorMessage(int(e.Status)))
}

// HTTPErrorMessages maps status codes to user-friendly messages
var HTTPErrorMessages = map[int]string{
	http.StatusBadRequest:          "Invalid request. ",
	http.StatusUnauthorized:        "Unauthorized. Invalid credentials.",
	http.StatusForbidden:           "Access denied. Check permissions.",
	http.StatusNotFound:            "Resource not found.",
	http.StatusInternalServerError: "Internal Server side error.",
	http.StatusBadGateway:          "Bad Gateway. ",
	http.StatusServiceUnavailable:  "Service is temporarily unavailable.",
}

// GetErrorMessage returns a message based on the HTTP status code
func getErrorMessage(statusCode int) string {
	if msg, exists := HTTPErrorMessages[statusCode]; exists {
		return msg
	}
	return "An unknown error occurred."
}
