package usecase

import (
	"context"
	"password-validator/core/domain/password"
	_errors "password-validator/core/errors"
	"password-validator/core/repository"
	"password-validator/core/usecase/input"
	"password-validator/core/usecase/output"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	testTable struct {
		name       string
		in         any
		repoReturn any
		repoErr    error
		out        any
		err        error
	}

	validatePasswordPresenterMock struct {
		mock.Mock
	}
)

func (u *validatePasswordPresenterMock) Output(ctx context.Context, p *password.Password) output.PasswordOutput {
	return output.PasswordOutput{
		IsValid: p.IsValid(),
	}
}

func TestValidatePasswordUseCase(t *testing.T) {
	tt := []testTable{
		{
			name: "successfully password creation",
			in: input.PasswordInput{
				Password: "AbTp9!fok",
			},
			repoErr: nil,
			out: output.PasswordOutput{
				IsValid: true,
			},
			err: nil,
		},
		{
			name: "validation error",
			in: input.PasswordInput{
				Password: "AbTp9!foA",
			},
			repoErr: nil,
			out:     output.PasswordOutput{},
			err:     _errors.InvalidField{Field: "password", AsIs: "Must not contain repeated characters (excluding spaces)"},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			repo := &repository.PasswordRepositoryMock{}
			repo.On("Save", mock.Anything, mock.Anything).Return(test.repoErr)
			uc := NewValidatePasswordUseCase(10*time.Second, repo, &validatePasswordPresenterMock{})
			out, err := uc.Execute(context.Background(), test.in.(input.PasswordInput))
			if test.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.out.(output.PasswordOutput).IsValid, out.IsValid)
			} else {
				assert.Equal(t, test.err, err)
			}
		})
	}
}
