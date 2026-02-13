package repository

import (
	"context"
	"password-validator/core/domain/password"
	_errors "password-validator/core/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveRepository(t *testing.T) {
	repo := NewPasswordRepository()
	p, _ := password.New(password.WithPassword("123"))
	repo.Save(context.TODO(), p)

	assert.Equal(t, repo.storage[0].IsValid(), p.IsValid())
}

func TestFindByIdRepository(t *testing.T) {
	p, _ := password.New(password.WithPassword("123"))
	tt := []struct {
		name        string
		input       string
		output      *password.Password
		expectedErr error
	}{
		{
			name:        "succes execution",
			input:       "123",
			output:      p,
			expectedErr: nil,
		},
		{
			name:   "password not found",
			input:  "100",
			output: &password.Password{},
			expectedErr: _errors.NotFoundError{
				Entity: "Password",
				ID:     "100",
			},
		},
	}

	repo := NewPasswordRepository()
	repo.Save(context.TODO(), p)
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			ex, err := repo.FindById(context.TODO(), test.input)

			assert.Equal(t, err, test.expectedErr)
			assert.Equal(t, ex.IsValid(), test.output.IsValid())
		})
	}
}
