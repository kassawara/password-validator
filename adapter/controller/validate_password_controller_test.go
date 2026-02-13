package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	_errors "password-validator/core/errors"
	"password-validator/core/usecase/input"
	"password-validator/core/usecase/output"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	ValidatePasswordUseCaseMock struct {
		mock.Mock
	}
	ErrReader int
)

func (c *ValidatePasswordUseCaseMock) Execute(ctx context.Context, i input.PasswordInput) (output.PasswordOutput, error) {
	ret := c.Called(ctx, i)
	return ret.Get(0).(output.PasswordOutput), ret.Error(1)
}

func (ErrReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func (ErrReader) Close() error {
	return nil
}

func TestValidatePasswordController(t *testing.T) {
	tt := []struct {
		name               string
		usecaseOutput      output.PasswordOutput
		stringBody         string
		expectedReadAllErr bool
		usecaseError       error
		expectedStatus     int
	}{
		{
			name:               "valid password",
			expectedReadAllErr: false,
			stringBody:         `{"password":"AbTp9!fok"}`,
			usecaseOutput: output.PasswordOutput{
				IsValid: true,
			},
			usecaseError:   nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:               "invalid password (less than 9 chars)",
			expectedReadAllErr: false,
			stringBody:         "{\"password\":\"12345678\"}",
			usecaseOutput:      output.PasswordOutput{},
			usecaseError:       _errors.InvalidField{Field: "password", AsIs: "Must have at least 9 characters"},
			expectedStatus:     http.StatusUnprocessableEntity,
		},
		{
			name:               "reading request body error",
			expectedReadAllErr: true,
			stringBody:         "",
			usecaseOutput:      output.PasswordOutput{},
			usecaseError:       nil,
			expectedStatus:     http.StatusBadRequest,
		},
		{
			name:               "unmarshal error",
			expectedReadAllErr: false,
			stringBody:         `error`,
			usecaseOutput:      output.PasswordOutput{},
			usecaseError:       errors.New("test"),
			expectedStatus:     http.StatusInternalServerError,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			var body io.ReadCloser
			stringReader := strings.NewReader(test.stringBody)
			if test.expectedReadAllErr {
				body = ErrReader(0)
			} else {
				body = io.NopCloser(stringReader)
			}
			req := &http.Request{
				Header: http.Header{},
				Body:   body,
			}
			uc := &ValidatePasswordUseCaseMock{}
			uc.On("Execute", mock.Anything, mock.Anything).Return(test.usecaseOutput, test.usecaseError)
			c := NewValidatePasswordController(uc)

			c.Execute(w, req)

			assert.Equal(t, w.Result().StatusCode, test.expectedStatus)
		})
	}
}
