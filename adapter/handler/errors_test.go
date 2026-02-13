package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	_errors "password-validator/core/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleErrors(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		statusCode int
	}{
		{
			name:       "NotFoundError should return status not found",
			err:        _errors.NotFoundError{},
			statusCode: http.StatusNotFound,
		},
		{
			name:       "InvalidField should return status bad request",
			err:        _errors.InvalidField{},
			statusCode: http.StatusUnprocessableEntity,
		},
		{
			name:       "Other error should return status internal server error",
			err:        errors.New("other error"),
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleErrors(w, test.err, nil)

			assert.Equal(t, w.Result().StatusCode, test.statusCode)
		})
	}

}
