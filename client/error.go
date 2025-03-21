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
	return fmt.Sprintf("HttpError: %d: %s", e.Status, getErrorMessage(int(e.Status)))
}

// HTTPErrorMessages maps status codes to user-friendly messages
var HTTPErrorMessages = map[int]string{
	http.StatusBadRequest:          "Invalid request. Please check your input.",
	http.StatusUnauthorized:        "Unauthorized. Please provide valid credentials.",
	http.StatusForbidden:           "Access denied. You do not have permission.",
	http.StatusNotFound:            "Resource not found.",
	http.StatusInternalServerError: "Internal Server side error. Please try again later.",
	http.StatusBadGateway:          "Bad Gateway. Please try again later.",
	http.StatusServiceUnavailable:  "Service is temporarily unavailable. Please try again later.",
}

// GetErrorMessage returns a message based on the HTTP status code
func getErrorMessage(statusCode int) string {
	if msg, exists := HTTPErrorMessages[statusCode]; exists {
		return msg
	}
	return "An unknown error occurred."
}
