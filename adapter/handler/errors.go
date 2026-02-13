package handler

import (
	"errors"
	"net/http"
	"password-validator/adapter/response"
	_errors "password-validator/core/errors"
)

func HandleErrors(w http.ResponseWriter, err error, output interface{}) {
	var status int
	switch {
	case errors.As(err, &_errors.NotFoundError{}):
		status = http.StatusNotFound
	case errors.As(err, &_errors.InvalidField{}):
		status = http.StatusUnprocessableEntity
	default:
		status = http.StatusInternalServerError
	}
	response.NewError(err, status, output).Send(w)
}
