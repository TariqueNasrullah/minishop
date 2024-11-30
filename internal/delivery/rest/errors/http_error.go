package errors

import "net/http"

type HTTPError struct {
	Message string              `json:"message,omitempty"`
	Type    string              `json:"type,omitempty"`
	Code    int                 `json:"code,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	Unauthrized = HTTPError{Message: "Unauthorized", Type: "error", Code: http.StatusUnauthorized}
)
