package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	statusCode int
	Error      string      `json:"error"`
	Output     interface{} `json:"password,omitempty"`
}

func NewError(err error, status int, output interface{}) *Error {
	return &Error{
		Error:      err.Error(),
		statusCode: status,
		Output:     output,
	}
}

func (err Error) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.statusCode)
	json.NewEncoder(w).Encode(Error{
		Error:  err.Error,
		Output: err.Output,
	})
}
