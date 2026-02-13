package presenter

import (
	"context"
	"password-validator/core/domain/password"
	"password-validator/core/usecase/output"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePasswordPresenter(t *testing.T) {
	p, _ := password.New(password.WithPassword("123"))
	tt := []struct {
		name   string
		input  *password.Password
		output output.PasswordOutput
	}{
		{
			name:   "success parse",
			input:  p,
			output: output.PasswordOutput{IsValid: p.IsValid()},
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			p := NewValidatePasswordPresenter()
			out := p.Output(context.TODO(), test.input)

			assert.Equal(t, test.output.IsValid, out.IsValid)
		})
	}
}
